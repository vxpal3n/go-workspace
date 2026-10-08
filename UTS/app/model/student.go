package model

import "time"

type Student struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	NIM         string     `json:"nim"`
	Nama        string     `json:"nama"`
	Prodi       string     `json:"prodi"`
	Angkatan    int        `json:"angkatan"`
	IPKTerakhir float64    `json:"ipk_terakhir"`
	DeletedAt   *time.Time `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateStudentRequest struct {
	NIM         string  `json:"nim"          validate:"required,len=12,numeric"`
	Nama        string  `json:"nama"         validate:"required,min=2,max=100"`
	Email       string  `json:"email"        validate:"required,email,max=255"`
	Prodi       string  `json:"prodi"        validate:"required,min=2,max=100"`
	Angkatan    int     `json:"angkatan"     validate:"required,min=2000,max=2100"`
	IPKTerakhir float64 `json:"ipk_terakhir" validate:"omitempty,min=0,max=4"`
}

type UpdateStudentRequest struct {
	Nama        string  `json:"nama"         validate:"required,min=2,max=100"`
	Prodi       string  `json:"prodi"        validate:"required,min=2,max=100"`
	Angkatan    int     `json:"angkatan"     validate:"required,min=2000,max=2100"`
	IPKTerakhir float64 `json:"ipk_terakhir" validate:"omitempty,min=0,max=4"`
}

type ListStudentsQuery struct {
	Page     int
	PerPage  int
	Prodi    string
	Angkatan int
	Search   string
	Sort     string
}

type StudentDetail struct {
	ID          int              `json:"id"`
	NIM         string           `json:"nim"`
	Nama        string           `json:"nama"`
	Prodi       string           `json:"prodi"`
	Angkatan    int              `json:"angkatan"`
	IPKTerakhir float64          `json:"ipk_terakhir"`
	MataKuliah  []EnrolledCourse `json:"mata_kuliah"`
	TotalSKS    int              `json:"total_sks"`
	BatasSKS    int              `json:"batas_sks"`
}

type EnrolledCourse struct {
	EnrollmentID  int    `json:"enrollment_id"`
	CourseID      int    `json:"course_id"`
	KodeMK        string `json:"kode_mk"`
	NamaMK        string `json:"nama_mk"`
	SKS           int    `json:"sks"`
	Semester      int    `json:"semester"`
	TahunAkademik string `json:"tahun_akademik"`
}