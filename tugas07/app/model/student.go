package model

import "time"

type Student struct {
	ID        int       `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Grade     float64   `json:"grade"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	OwnerID   *int      `json:"owner_id,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateStudentRequest struct {
	NIM      string  `json:"nim"      validate:"required,min=3,max=20,alphanum"`
	Name     string  `json:"name"     validate:"required,min=2,max=100"`
	Email    string  `json:"email"    validate:"required,email,max=120"`
	Grade    float64 `json:"grade"    validate:"min=0,max=100"`
	Password string  `json:"password" validate:"required,min=8,max=72,nospace"`
}

type ReplaceStudentRequest struct {
	NIM      string  `json:"nim"      validate:"required,min=3,max=20,alphanum"`
	Name     string  `json:"name"     validate:"required,min=2,max=100"`
	Email    string  `json:"email"    validate:"required,email,max=120"`
	Grade    float64 `json:"grade"    validate:"min=0,max=100"`
	IsActive bool    `json:"is_active"`
}

type PatchStudentRequest struct {
	NIM      *string  `json:"nim,omitempty"      validate:"omitnil,min=3,max=20,alphanum"`
	Name     *string  `json:"name,omitempty"     validate:"omitnil,min=2,max=100"`
	Email    *string  `json:"email,omitempty"    validate:"omitnil,email,max=120"`
	Grade    *float64 `json:"grade,omitempty"    validate:"omitnil,min=0,max=100"`
	IsActive *bool    `json:"is_active,omitempty"`
}

type AssignRoleRequest struct {
	Role string `json:"role"`
}

type WebResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ListQuery struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	IsActive *bool
	MinGrade *float64
	MaxGrade *float64
}

type ErrorResponse struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

type Cursor struct {
	CreatedAt time.Time
	ID        int
}

type CursorQuery struct {
	Search   string
	IsActive *bool
	After    *Cursor
	Limit    int
}

type CursorMeta struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

func (q ListQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}