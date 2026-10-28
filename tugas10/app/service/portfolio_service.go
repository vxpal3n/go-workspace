package service

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"tugas10/app/model"
	"tugas10/app/repository"
	"tugas10/helper"
)

type PortfolioService struct {
	portfolios repository.PortfolioRepository
	students   repository.StudentRepository
	perms      *helper.PermissionSet
}

func NewPortfolioService(
	portfolios repository.PortfolioRepository,
	students repository.StudentRepository,
	perms *helper.PermissionSet,
) *PortfolioService {
	return &PortfolioService{portfolios: portfolios, students: students, perms: perms}
}

func (s *PortfolioService) Get(c *fiber.Ctx) error {
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

	if !CanAccessStudent(current, id, nil, s.perms, "student:read:any") {
		return helper.Forbidden("tidak berhak mengakses portofolio mahasiswa lain")
	}

	if _, err := s.students.FindByID(ctx, id); err != nil {
		return translatePortfolioError(err)
	}

	portfolio, err := s.portfolios.FindByStudent(ctx, id)
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.Success(c, fiber.StatusOK, "portofolio ditemukan", portfolio)
}

func (s *PortfolioService) Put(c *fiber.Ctx) error {
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

	if !CanAccessStudent(current, id, nil, s.perms, "student:update:any") {
		return helper.Forbidden("tidak berhak mengubah portofolio mahasiswa lain")
	}

	if _, err := s.students.FindByID(ctx, id); err != nil {
		return translatePortfolioError(err)
	}

	var req model.PutPortfolioRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	req.Headline = strings.TrimSpace(req.Headline)
	req.Bio = strings.TrimSpace(req.Bio)

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	portfolio := model.Portfolio{
		StudentID: id,
		Headline:  req.Headline,
		Bio:       req.Bio,
		Projects:  make([]model.Project, 0, len(req.Projects)),
		Certificates: make([]model.Certificate, 0, len(req.Certificates)),
	}
	for _, p := range req.Projects {
		techs := NormalizeTags(p.Technologies)
		portfolio.Projects = append(portfolio.Projects, model.Project{
			Title: p.Title, Technologies: techs, Link: p.Link, Year: p.Year,
		})
	}
	for _, cert := range req.Certificates {
		portfolio.Certificates = append(portfolio.Certificates, model.Certificate{
			Name: cert.Name, Issuer: cert.Issuer, Date: cert.Date,
		})
	}

	saved, err := s.portfolios.Upsert(ctx, portfolio)
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.Success(c, fiber.StatusOK, "portofolio berhasil disimpan", saved)
}

func (s *PortfolioService) AddProject(c *fiber.Ctx) error {
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

	if !CanAccessStudent(current, id, nil, s.perms, "student:update:any") {
		return helper.Forbidden("tidak berhak menambah proyek mahasiswa lain")
	}

	if _, err := s.students.FindByID(ctx, id); err != nil {
		return translatePortfolioError(err)
	}

	var req model.AddProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	techs := NormalizeTags(req.Technologies)
	updated, err := s.portfolios.AddProject(ctx, id, model.Project{
		Title: req.Title, Technologies: techs, Link: req.Link, Year: req.Year,
	})
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.Success(c, fiber.StatusCreated, "proyek berhasil ditambahkan", updated)
}

func (s *PortfolioService) DeleteProject(c *fiber.Ctx) error {
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
	if !CanAccessStudent(current, id, nil, s.perms, "student:update:any") {
		return helper.Forbidden("tidak berhak mengubah portofolio mahasiswa lain")
	}

	projectID := strings.TrimSpace(c.Params("projectId"))
	if projectID == "" {
		return helper.BadRequest("projectId wajib diisi")
	}

	updated, err := s.portfolios.DeleteProject(ctx, id, projectID)
	if err != nil {
		return translatePortfolioError(err)
	}
	return helper.Success(c, fiber.StatusOK, "proyek berhasil dihapus", updated)
}

func (s *PortfolioService) StatsTechnologies(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	limit := c.QueryInt("limit", 10)
	stats, err := s.portfolios.StatsTechnologies(ctx, limit)
	if err != nil {
		return helper.Internal(err)
	}
	return helper.Success(c, fiber.StatusOK, "teknologi terpopuler berhasil diambil", stats)
}

func translatePortfolioError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return helper.NotFound("mahasiswa tidak ditemukan")
	case errors.Is(err, repository.ErrPortfolioNotFound):
		return helper.NotFound("portofolio belum dibuat untuk mahasiswa ini")
	case errors.Is(err, repository.ErrProjectLimit):
		return helper.NewError(fiber.StatusConflict, "PROJECT_LIMIT",
			"jumlah proyek sudah mencapai batas (20)")
	case errors.Is(err, repository.ErrInvalidID):
		return helper.BadRequest("id proyek tidak valid")
	default:
		return helper.Internal(err)
	}
}