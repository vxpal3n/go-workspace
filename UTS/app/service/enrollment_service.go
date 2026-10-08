package service

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"

	"UTS/app/model"
	"UTS/app/repository"
	"UTS/helper"
)

type EnrollmentService struct {
	students    repository.StudentRepository
	courses     repository.CourseRepository
	enrollments repository.EnrollmentRepository
}

func NewEnrollmentService(
	students repository.StudentRepository,
	courses repository.CourseRepository,
	enrollments repository.EnrollmentRepository,
) *EnrollmentService {
	return &EnrollmentService{
		students: students, courses: courses, enrollments: enrollments,
	}
}

func (s *EnrollmentService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}
	if current.Role != "mahasiswa" {
		return helper.Forbidden("hanya mahasiswa yang dapat mengambil mata kuliah")
	}

	var req model.CreateEnrollmentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	req.TahunAkademik = strings.TrimSpace(req.TahunAkademik)

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}
	if err := ParseTahunAkademik(req.TahunAkademik); err != nil {
		return helper.Validation(map[string]string{
			"tahun_akademik": err.Error(),
		})
	}

	student, err := s.students.FindByUserID(ctx, current.UserID)
	if err != nil {
		return helper.Forbidden("Data mahasiswa tidak ditemukan")
	}

	tx, err := s.enrollments.Begin(ctx)
	if err != nil {
		return helper.Internal(err)
	}
	defer tx.Rollback(ctx)

	course, err := s.enrollments.LockCourseForUpdate(ctx, tx, req.CourseID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Validation(map[string]string{
				"course_id": "Mata kuliah tidak ditemukan",
			})
		}
		return helper.Internal(err)
	}

	exists, err := s.enrollments.Exists(ctx, tx, student.ID, course.ID, req.TahunAkademik)
	if err != nil {
		return helper.Internal(err)
	}
	if exists {
		return helper.Conflict("Anda sudah mengambil mata kuliah ini pada tahun akademik yang sama")
	}

	terisi, err := s.enrollments.CountByCourseAndYear(ctx, tx, course.ID, req.TahunAkademik)
	if err != nil {
		return helper.Internal(err)
	}
	if terisi >= course.Kuota {
		return helper.Validation(map[string]string{
			"course_id": "Kuota mata kuliah sudah penuh",
		})
	}

	currentSKS, err := s.enrollments.SumSKSForStudentYear(ctx, tx, student.ID, req.TahunAkademik)
	if err != nil {
		return helper.Internal(err)
	}
	batas := BatasSKS(student.IPKTerakhir)
	if err := CanEnrollSKS(currentSKS, course.SKS, batas); err != nil {
		return helper.Validation(map[string]string{
			"course_id": err.Error(),
		})
	}

	enrollment, err := s.enrollments.Insert(ctx, tx, student.ID, course.ID, req.TahunAkademik)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("Anda sudah mengambil mata kuliah ini pada tahun akademik yang sama")
		}
		return helper.Internal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return helper.Internal(err)
	}

	return helper.Created(c, "Mata kuliah berhasil diambil", model.EnrollmentDetail{
		ID:            enrollment.ID,
		CourseID:      course.ID,
		KodeMK:        course.KodeMK,
		NamaMK:        course.NamaMK,
		SKS:           course.SKS,
		TahunAkademik: enrollment.TahunAkademik,
		CreatedAt:     enrollment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (s *EnrollmentService) Delete(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	current, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}
	if current.Role != "mahasiswa" {
		return helper.Forbidden("hanya mahasiswa yang dapat membatalkan KRS")
	}

	id, valid := helper.ParamID(c)
	if !valid {
		return helper.BadRequest("id harus berupa angka positif")
	}

	enrollment, err := s.enrollments.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("Enrollment tidak ditemukan")
		}
		return helper.Internal(err)
	}

	student, err := s.students.FindByUserID(ctx, current.UserID)
	if err != nil || student.ID != enrollment.StudentID {
		return helper.Forbidden("Anda tidak berhak membatalkan enrollment ini")
	}

	if err := s.enrollments.Delete(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, repository.ErrNotFound) {
			return helper.NotFound("Enrollment tidak ditemukan")
		}
		return helper.Internal(err)
	}

	return helper.NoContent(c)
}