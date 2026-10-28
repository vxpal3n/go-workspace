package model

import "time"

type Portfolio struct {
	ID           string        `json:"id"`
	StudentID    int           `json:"student_id"`
	Headline     string        `json:"headline"`
	Bio          string        `json:"bio"`
	Projects     []Project     `json:"projects"`
	Certificates []Certificate `json:"certificates"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type Project struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Technologies []string `json:"technologies"`
	Link         string   `json:"link,omitempty"`
	Year         int      `json:"year"`
}

type Certificate struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Issuer string `json:"issuer"`
	Date   string `json:"date"`
}

type PutPortfolioRequest struct {
	Headline     string               `json:"headline"     validate:"max=120"`
	Bio          string               `json:"bio"          validate:"max=1000"`
	Projects     []ProjectRequest     `json:"projects"     validate:"max=20,dive"`
	Certificates []CertificateRequest `json:"certificates" validate:"max=20,dive"`
}

type ProjectRequest struct {
	Title        string   `json:"title"        validate:"required,min=2,max=120"`
	Technologies []string `json:"technologies" validate:"max=10,dive,min=1,max=30"`
	Link         string   `json:"link"         validate:"omitempty,url,max=255"`
	Year         int      `json:"year"         validate:"required,min=1900,max=2100"`
}

type CertificateRequest struct {
	Name   string `json:"name"   validate:"required,min=2,max=120"`
	Issuer string `json:"issuer" validate:"required,min=2,max=120"`
	Date   string `json:"date"   validate:"required,datetime=2006-01-02"`
}

type AddProjectRequest struct {
	Title        string   `json:"title"        validate:"required,min=2,max=120"`
	Technologies []string `json:"technologies" validate:"max=10,dive,min=1,max=30"`
	Link         string   `json:"link"         validate:"omitempty,url,max=255"`
	Year         int      `json:"year"         validate:"required,min=1900,max=2100"`
}

type TechCount struct {
	Technology string `json:"technology"`
	Total      int    `json:"total"`
}