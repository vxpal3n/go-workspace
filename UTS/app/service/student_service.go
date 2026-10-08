package service

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"UTS/app/model"
	"UTS/app/repository"
	"UTS/helper"
)

type StudentService struct {
	students    repository.StudentRepository
	enrollments repository.EnrollmentRepository
}

func NewStudentService(
	students repository.StudentRepository,
	enrollments repository.EnrollmentRepository,
) *StudentService {
	return &StudentService{students: students, enrollments: enrollments}
}

func translateError(err error, entity string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound(entity + " tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Validation(map[string]string{
			"nim":   "NIM sudah terdaftar",
			"email": "Email sudah terdaftar",
		})
	case errors.Is(err, repository.ErrConflict):
		return helper.Conflict("Operasi bertentangan dengan data yang ada")
	default:
		return helper.Internal(err)
	}
}

func (s *StudentService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	q := model.ListStudentsQuery{
		Page:     c.QueryInt("page", 1),
		PerPage:  c.QueryInt("per_page", 10),
		Prodi:    strings.TrimSpace(c.Query("prodi")),
		Angkatan: c.QueryInt("angkatan", 0),
		Search:   strings.TrimSpace(c.Query("search")),
		Sort:     c.Query("sort", "nama"),
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.PerPage < 1 {
		q.PerPage = 10
	}
	if q.PerPage > 50 {
		q.PerPage = 50
	}

	students, total, err := s.students.List(ctx, q)
	if err != nil {
		return helper.Internal(err)
	}

	lastPage := (total + q.PerPage - 1) / q.PerPage
	if lastPage == 0 {
		lastPage = 1
	}

	return helper.SuccessList(c, "Data mahasiswa berhasil diambil", students, &model.Meta{
		CurrentPage: q.Page,
		PerPage:     q.PerPage,
		Total:       total,
		LastPage:    lastPage,
	})
}

func (s *StudentService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.NIM = strings.TrimSpace(req.NIM)
	req.Nama = strings.TrimSpace(req.Nama)
	req.Prodi = strings.TrimSpace(req.Prodi)

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	hashed, err := helper.HashPassword(req.NIM)
	if err != nil {
		return helper.Internal(err)
	}

	created, err := s.students.CreateWithUser(ctx, req.Email, hashed, model.Student{
		NIM:         req.NIM,
		Nama:        req.Nama,
		Prodi:       req.Prodi,
		Angkatan:    req.Angkatan,
		IPKTerakhir: req.IPKTerakhir,
	})
	if err != nil {
		return translateError(err, "Mahasiswa")
	}

	return helper.Created(c, "Mahasiswa berhasil ditambahkan", created)
}

func (s *StudentService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	student, err := s.students.FindByID(ctx, id)
	if err != nil {
		return translateError(err, "Mahasiswa")
	}

	if current.Role == "mahasiswa" {
		own, err := s.students.FindByUserID(ctx, current.UserID)
		if err != nil || own.ID != student.ID {
			return helper.Forbidden("Anda tidak berhak mengakses data mahasiswa lain")
		}
	}

	courses, err := s.enrollments.ListEnrichedByStudent(ctx, student.ID)
	if err != nil {
		return helper.Internal(err)
	}

	sksList := make([]int, len(courses))
	for i, c := range courses {
		sksList[i] = c.SKS
	}

	return helper.Success(c, fiber.StatusOK, "Detail mahasiswa berhasil diambil", model.StudentDetail{
		ID:          student.ID,
		NIM:         student.NIM,
		Nama:        student.Nama,
		Prodi:       student.Prodi,
		Angkatan:    student.Angkatan,
		IPKTerakhir: student.IPKTerakhir,
		MataKuliah:  courses,
		TotalSKS:    TotalSKS(sksList),
		BatasSKS:    BatasSKS(student.IPKTerakhir),
	})
}

func (s *StudentService) Update(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	var req model.UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	req.Nama = strings.TrimSpace(req.Nama)
	req.Prodi = strings.TrimSpace(req.Prodi)

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	updated, err := s.students.Update(ctx, model.Student{
		ID:          id,
		Nama:        req.Nama,
		Prodi:       req.Prodi,
		Angkatan:    req.Angkatan,
		IPKTerakhir: req.IPKTerakhir,
	})
	if err != nil {
		return translateError(err, "Mahasiswa")
	}

	return helper.Success(c, fiber.StatusOK, "Data mahasiswa berhasil diperbarui", updated)
}

func (s *StudentService) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	if err := s.students.SoftDelete(ctx, id); err != nil {
		return translateError(err, "Mahasiswa")
	}

	return helper.NoContent(c)
}