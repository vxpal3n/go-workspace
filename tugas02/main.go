package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var methodsWithBody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

func requireJSON(c *fiber.Ctx) error {
	if methodsWithBody[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return fail(c, fiber.StatusUnsupportedMediaType, "Content-Type harus application/json")
		}
	}
	return c.Next()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "API Students - Praktikum Backend Lanjut",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			message := "terjadi kesalahan pada server"
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}
			return fail(c, code, message)
		},
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${time}] ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
	}))
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! from API Students")
	})

	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return ok(c, "server berjalan", fiber.Map{"timestamp": time.Now()})
	})

	studentsGroup := api.Group("/students", requireJSON)
	studentsGroup.Get("/", listStudents)
	studentsGroup.Get("/:id", getStudent)
	studentsGroup.Post("/", createStudent)
	studentsGroup.Put("/:id", replaceStudent)
	studentsGroup.Patch("/:id", patchStudent)
	studentsGroup.Delete("/:id", deleteStudent)

	app.Use(func(c *fiber.Ctx) error {
		return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	fmt.Println("Server berjalan di http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}