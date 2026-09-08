package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"tugas03/config"
	"tugas03/database"
	"tugas03/app/repository"
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
	config.LoadEnv()

	ctx := context.Background()
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	studentRepo := repository.NewStudentRepository(pool)
	studentHandler := NewStudentHandler(studentRepo)

	app := fiber.New(fiber.Config{
		AppName: "API Students - Modul 3 (PostgreSQL)",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := "terjadi kesalahan pada server"
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				msg = e.Message
			}
			return fail(c, code, msg)
		},
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${time}] ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
	}))
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! from API Students (Modul 3)")
	})

	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			return fail(c, fiber.StatusServiceUnavailable, "database tidak dapat dihubungi")
		}
		return ok(c, "server dan database berjalan", map[string]interface{}{"timestamp": time.Now()})
	})

	students := api.Group("/students", requireJSON)
	students.Get("/", studentHandler.List)
	students.Get("/:id", studentHandler.Get)
	students.Post("/", studentHandler.Create)
	students.Put("/:id", studentHandler.Replace)
	students.Patch("/:id", studentHandler.Patch)
	students.Delete("/:id", studentHandler.Delete)

	app.Use(func(c *fiber.Ctx) error {
		return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
	})

	port := config.GetEnv("APP_PORT", "3000")
	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}