package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"tugas03/app/model"
)


func ok(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func okList(c *fiber.Ctx, message string, data interface{}, meta *model.Meta) error {
	return c.Status(fiber.StatusOK).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func created(c *fiber.Ctx, message string, data interface{}, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(model.WebResponse{
		Success: false,
		Message: message,
	})
}

func failValidation(c *fiber.Ctx, errors map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(model.WebResponse{
		Success: false,
		Message: "validasi gagal",
		Errors:  errors,
	})
}


func parseListQuery(c *fiber.Ctx) model.ListQuery {
	q := model.ListQuery{
		Page:  1,
		Limit: 10,
		Sort:  "id",
		Order: "asc",
	}
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			q.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			if l > 100 {
				l = 100
			}
			q.Limit = l
		}
	}
	if search := c.Query("search"); search != "" {
		q.Search = strings.TrimSpace(search)
	}
	if sort := c.Query("sort"); sort != "" {
		switch sort {
		case "id", "nim", "name", "grade", "created_at":
			q.Sort = sort
		}
	}
	if order := c.Query("order"); order != "" {
		if order == "asc" || order == "desc" {
			q.Order = order
		}
	}
	if isActive := c.Query("is_active"); isActive != "" {
		if b, err := strconv.ParseBool(isActive); err == nil {
			q.IsActive = &b
		}
	}
	if minGrade := c.Query("min_grade"); minGrade != "" {
		if g, err := strconv.ParseFloat(minGrade, 64); err == nil {
			q.MinGrade = &g
		}
	}
	if maxGrade := c.Query("max_grade"); maxGrade != "" {
		if g, err := strconv.ParseFloat(maxGrade, 64); err == nil {
			q.MaxGrade = &g
		}
	}
	return q
}


func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}


func reqCtx(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}