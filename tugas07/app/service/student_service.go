package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"tugas07/app/model"
	"tugas07/app/repository"
	"tugas07/helper"
)

type StudentService struct {
	repo  repository.StudentRepository
	perms *helper.PermissionSet
}

func NewStudentService(repo repository.StudentRepository, perms *helper.PermissionSet) *StudentService {
	return &StudentService{repo: repo, perms: perms}
}

func translateError(err error, entity string) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound(entity + " tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return helper.Conflict("NIM atau email sudah dipakai")
	default:
		return nil
	}
}

func (s *StudentService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	q := helper.ParseListQuery(c)
	students, total, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return helper.Internal(err)
	}
	return helper.SuccessList(c, "daftar student berhasil diambil", students, &model.Meta{
		Page:       q.Page,
		Limit:      q.Limit,
		Total:      total,
		TotalPages: CountTotalPages(total, q.Limit),
	})
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

	target, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(err, "student")
	}

	if !CanAccessStudent(current, target.ID, target.OwnerID, s.perms, "student:read:any") {
		return helper.Forbidden("tidak berhak mengakses data student ini")
	}
	return helper.Success(c, fiber.StatusOK, "student ditemukan", target)
}

func (s *StudentService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if errs := ValidateCreateByAdmin(req); len(errs) > 0 {
		return helper.Validation(errs)
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.Internal(err)
	}

	created, err := s.repo.CreateWithOwner(ctx, model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Email:    req.Email,
		Grade:    req.Grade,
		Password: hashed,
		Role:     "user",
		IsActive: true,
	}, current.StudentID)
	if err != nil {
		return translateError(err, "student")
	}

	location := c.BaseURL() + "/api/v1/students/" + strconv.Itoa(created.ID)
	return helper.Created(c, "student berhasil dibuat", created, location)
}

func (s *StudentService) Replace(c *fiber.Ctx) error {
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

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	errs := map[string]string{}
	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if !isValidEmail(req.Email) {
		errs["email"] = "format email tidak valid"
	}
	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "harus antara 0 dan 100"
	}
	if len(errs) > 0 {
		return helper.Validation(errs)
	}

	target, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(err, "student")
	}
	if !CanAccessStudent(current, target.ID, target.OwnerID, s.perms, "student:update:any") {
		return helper.Forbidden("tidak berhak mengubah data student ini")
	}

	updated := model.Student{
		ID:       id,
		NIM:      req.NIM,
		Name:     req.Name,
		Email:    req.Email,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	}
	result, err := s.repo.Update(ctx, updated)
	if err != nil {
		return translateError(err, "student")
	}
	return helper.Success(c, fiber.StatusOK, "student berhasil diganti seluruhnya", result)
}

func (s *StudentService) Patch(c *fiber.Ctx) error {
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

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if req.NIM == nil && req.Name == nil && req.Email == nil && req.Grade == nil && req.IsActive == nil {
		return helper.BadRequest("tidak ada field yang diubah")
	}

	target, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return translateError(err, "student")
	}
	if !CanAccessStudent(current, target.ID, target.OwnerID, s.perms, "student:update:any") {
		return helper.Forbidden("tidak berhak mengubah data student ini")
	}

	errs := map[string]string{}
	if req.NIM != nil {
		nim := strings.TrimSpace(*req.NIM)
		if nim == "" {
			errs["nim"] = "tidak boleh kosong"
		} else {
			target.NIM = nim
		}
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			errs["name"] = "tidak boleh kosong"
		} else {
			target.Name = name
		}
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if !isValidEmail(email) {
			errs["email"] = "format email tidak valid"
		} else {
			target.Email = email
		}
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 100 {
			errs["grade"] = "harus antara 0 dan 100"
		} else {
			target.Grade = *req.Grade
		}
	}
	if req.IsActive != nil {
		target.IsActive = *req.IsActive
	}
	if len(errs) > 0 {
		return helper.Validation(errs)
	}

	result, err := s.repo.Update(ctx, target)
	if err != nil {
		return translateError(err, "student")
	}
	return helper.Success(c, fiber.StatusOK, "student berhasil diperbarui sebagian", result)
}

func (s *StudentService) Delete(c *fiber.Ctx) error {
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

	if current.StudentID == id {
		return helper.Forbidden("tidak boleh menghapus akun sendiri")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return translateError(err, "student")
	}
	return helper.NoContent(c)
}

func (s *StudentService) AssignRole(c *fiber.Ctx) error {
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

	var req model.AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}

	if errs := ValidateAssignRole(current, id, req, s.perms); len(errs) > 0 {
		return helper.Validation(errs)
	}

	result, err := s.repo.UpdateRole(ctx, id, strings.TrimSpace(req.Role))
	if err != nil {
		return translateError(err, "student")
	}
	return helper.Success(c, fiber.StatusOK, "role student berhasil diubah", result)
}