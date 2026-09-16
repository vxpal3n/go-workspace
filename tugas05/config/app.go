package config

import (
    "log/slog"

    "github.com/gofiber/fiber/v2"
    "github.com/jackc/pgx/v5/pgxpool"
    "tugas05/app/service"
    "tugas05/helper"
    "tugas05/middleware"
    "tugas05/route"
)

func NewApp(logger *slog.Logger, pool *pgxpool.Pool, studentService *service.StudentService) *fiber.App {
    app := fiber.New(fiber.Config{
        AppName: GetEnv("APP_NAME", "API Students - Modul 4 (Clean Architecture)"),
        ErrorHandler: newErrorHandler(logger),
    })

    middleware.Register(app, logger)
    route.Register(app, pool, studentService)

    app.Use(func(c *fiber.Ctx) error {
        return helper.Fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
    })

    return app
}

func newErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
    return func(c *fiber.Ctx, err error) error {
        status := fiber.StatusInternalServerError
        message := "terjadi error pada server"
        if e, ok := err.(*fiber.Error); ok {
            status = e.Code
            message = e.Message
        }
        logger.Error("unhandled_error",
            slog.String("path", c.Path()),
            slog.Int("status", status),
            slog.String("error", err.Error()),
        )
        return helper.Fail(c, status, message)
    }
}