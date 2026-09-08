package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"tugas03/app/model"
	"tugas03/app/repository"
)

type StudentHandler struct {
	repo repository.StudentRepository
}

func NewStudentHandler(repo repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

func translateError(c *fiber.Ctx, err error, defaultMsg string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return fail(c, fiber.StatusNotFound, "student tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return fail(c, fiber.StatusConflict, "NIM sudah digunakan")
	default:
		return fail(c, fiber.StatusInternalServerError, defaultMsg)
	}
}


func (h *StudentHandler) List(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	q := parseListQuery(c)
	students, total, err := h.repo.FindAll(ctx, q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "gagal mengambil data student")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = (total + q.Limit - 1) / q.Limit
	}
	return okList(c, "daftar student berhasil diambil", students, &model.Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}


func (h *StudentHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	student, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(c, err, "gagal mengambil data student")
	}
	return ok(c, "student ditemukan", student)
}


func (h *StudentHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	// Validasi
	errs := make(map[string]string)
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
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

	newStudent := model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: true,
	}
	saved, err := h.repo.Create(ctx, newStudent)
	if err != nil {
		return translateError(c, err, "gagal menyimpan student")
	}
	location := c.BaseURL() + "/api/v1/students/" + strconv.Itoa(saved.ID)
	return created(c, "student berhasil dibuat", saved, location)
}


func (h *StudentHandler) Replace(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	errs := make(map[string]string)
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
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

	updated := model.Student{
		ID:       id,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	}
	result, err := h.repo.Update(ctx, updated)
	if err != nil {
		return translateError(c, err, "gagal memperbarui student")
	}
	return ok(c, "student berhasil diganti seluruhnya", result)
}


func (h *StudentHandler) Patch(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}
	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	current, err := h.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(c, err, "gagal mengambil data student")
	}

	errs := make(map[string]string)
	if req.NIM != nil {
		nim := strings.TrimSpace(*req.NIM)
		if nim == "" {
			errs["nim"] = "tidak boleh kosong"
		} else {
			current.NIM = nim
		}
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			errs["name"] = "tidak boleh kosong"
		} else {
			current.Name = name
		}
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			errs["grade"] = "harus antara 0 dan 100"
		} else {
			current.Grade = *req.Grade
		}
	}
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	result, err := h.repo.Update(ctx, current)
	if err != nil {
		return translateError(c, err, "gagal memperbarui student")
	}
	return ok(c, "student berhasil diperbarui sebagian", result)
}


func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}
	if err := h.repo.Delete(ctx, id); err != nil {
		return translateError(c, err, "gagal menghapus student")
	}
	return noContent(c)
}