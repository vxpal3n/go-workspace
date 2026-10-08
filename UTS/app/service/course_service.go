package service

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"UTS/app/model"
	"UTS/app/repository"
	"UTS/helper"
)

type CourseService struct {
	courses repository.CourseRepository
}

func NewCourseService(courses repository.CourseRepository) *CourseService {
	return &CourseService{courses: courses}
}

func (s *CourseService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	q := model.ListCoursesQuery{
		Semester:  c.QueryInt("semester", 0),
		Search:    strings.TrimSpace(c.Query("search")),
		Available: c.Query("available") == "true",
	}

	courses, err := s.courses.List(ctx, q)
	if err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, fiber.StatusOK, "Data mata kuliah berhasil diambil", courses)
}