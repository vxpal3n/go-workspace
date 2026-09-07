package main

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var students []Student
var nextID = 1


func findStudentIndex(id int) int {
	for i := range students {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

func isNIMExists(nim string, excludeID int) bool {
	for _, s := range students {
		if s.ID != excludeID && strings.EqualFold(s.NIM, nim) {
			return true
		}
	}
	return false
}

func studentMatchesSearch(s Student, search string) bool {
	search = strings.ToLower(search)
	return strings.Contains(strings.ToLower(s.Name), search) ||
		strings.Contains(strings.ToLower(s.NIM), search)
}


func listStudents(c *fiber.Ctx) error {
	q := parseListQuery(c)

	filtered := []Student{}
	for _, s := range students {
		if q.IsActive != nil && s.IsActive != *q.IsActive {
			continue
		}
		if q.MinGrade != nil && s.Grade < *q.MinGrade {
			continue
		}
		if q.MaxGrade != nil && s.Grade > *q.MaxGrade {
			continue
		}
		if q.Search != "" && !studentMatchesSearch(s, q.Search) {
			continue
		}
		filtered = append(filtered, s)
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		var less bool
		switch q.Sort {
		case "nim":
			less = filtered[i].NIM < filtered[j].NIM
		case "name":
			less = filtered[i].Name < filtered[j].Name
		case "grade":
			less = filtered[i].Grade < filtered[j].Grade
		case "created_at":
			less = filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		default:
			less = filtered[i].ID < filtered[j].ID
		}
		if q.Order == "desc" {
			return !less
		}
		return less
	})

	total := len(filtered)
	totalPages := (total + q.Limit - 1) / q.Limit
	start := (q.Page - 1) * q.Limit
	if start > total {
		start = total
	}
	end := start + q.Limit
	if end > total {
		end = total
	}
	pageData := filtered[start:end]

	meta := &Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
	return okList(c, "daftar student berhasil diambil", pageData, meta)
}


func getStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}
	return ok(c, "student ditemukan", students[idx])
}


func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := make(map[string]string)
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	} else if isNIMExists(req.NIM, -1) {
		errs["nim"] = "sudah digunakan"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "harus antara 0 dan 100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	student := Student{
		ID:        nextID,
		NIM:       req.NIM,
		Name:      req.Name,
		Grade:     req.Grade,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	students = append(students, student)
	nextID++

	location := c.BaseURL() + "/api/v1/students/" + strconv.Itoa(student.ID)
	return created(c, "student berhasil dibuat", student, location)
}

func replaceStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}

	var req ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := make(map[string]string)
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
	} else if isNIMExists(req.NIM, id) {
		errs["nim"] = "sudah digunakan oleh student lain"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "harus antara 0 dan 100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	students[idx].NIM = req.NIM
	students[idx].Name = req.Name
	students[idx].Grade = req.Grade
	students[idx].IsActive = req.IsActive

	return ok(c, "student berhasil diganti seluruhnya", students[idx])
}


func patchStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}

	var req PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	errs := make(map[string]string)

	if req.NIM != nil {
		nim := strings.TrimSpace(*req.NIM)
		if nim == "" {
			errs["nim"] = "tidak boleh kosong"
		} else if isNIMExists(nim, id) {
			errs["nim"] = "sudah digunakan oleh student lain"
		} else {
			students[idx].NIM = nim
		}
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			errs["name"] = "tidak boleh kosong"
		} else {
			students[idx].Name = name
		}
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			errs["grade"] = "harus antara 0 dan 100"
		} else {
			students[idx].Grade = *req.Grade
		}
	}
	if req.IsActive != nil {
		students[idx].IsActive = *req.IsActive
	}

	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	return ok(c, "student berhasil diperbarui sebagian", students[idx])
}

func deleteStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	idx := findStudentIndex(id)
	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	}
	students = append(students[:idx], students[idx+1:]...)
	return noContent(c)
}