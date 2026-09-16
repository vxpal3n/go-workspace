package helper

import (
    "context"
    "strconv"
    "strings"
    "time"

    "github.com/gofiber/fiber/v2"
    "tugas05/app/model"
)

var allowedSort = map[string]bool{
    "id": true, "nim": true, "name": true, "grade": true, "created_at": true,
}

func RequestContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
    return context.WithTimeout(c.UserContext(), 5*time.Second)
}

func ParamID(c *fiber.Ctx) (int, bool) {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil || id < 1 {
        return 0, false
    }
    return id, true
}

func ParseListQuery(c *fiber.Ctx) model.ListQuery {
    q := model.ListQuery{
        Page:   1,
        Limit:  10,
        Search: strings.TrimSpace(c.Query("search")),
        Sort:   c.Query("sort", "id"),
        Order:  strings.ToLower(c.Query("order", "asc")),
    }
    if page := c.QueryInt("page", 1); page > 0 {
        q.Page = page
    }
    if limit := c.QueryInt("limit", 10); limit > 0 {
        if limit > 100 {
            limit = 100
        }
        q.Limit = limit
    }
    if !allowedSort[q.Sort] {
        q.Sort = "id"
    }
    if q.Order != "desc" {
        q.Order = "asc"
    }
    if raw := c.Query("is_active"); raw != "" {
        if v, err := strconv.ParseBool(raw); err == nil {
            q.IsActive = &v
        }
    }
    if raw := c.Query("min_grade"); raw != "" {
        if v, err := strconv.ParseFloat(raw, 64); err == nil {
            q.MinGrade = &v
        }
    }
    if raw := c.Query("max_grade"); raw != "" {
        if v, err := strconv.ParseFloat(raw, 64); err == nil {
            q.MaxGrade = &v
        }
    }
    return q
}