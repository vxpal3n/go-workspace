package helper

import (
	"testing"

	"UTS/app/model"
)

func TestValidateStruct_LoginRequest(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		req := model.LoginRequest{
			Email:    "admin@siakad.test",
			Password: "admin12345",
		}
		if errs := ValidateStruct(req); errs != nil {
			t.Errorf("harap lolos, dapat: %v", errs)
		}
	})

	t.Run("email salah format", func(t *testing.T) {
		req := model.LoginRequest{
			Email:    "bukan-email",
			Password: "admin12345",
		}
		errs := ValidateStruct(req)
		if _, ok := errs["email"]; !ok {
			t.Errorf("harap error pada email, dapat: %v", errs)
		}
	})

	t.Run("password terlalu pendek", func(t *testing.T) {
		req := model.LoginRequest{
			Email:    "admin@siakad.test",
			Password: "abc",
		}
		errs := ValidateStruct(req)
		if _, ok := errs["password"]; !ok {
			t.Errorf("harap error pada password, dapat: %v", errs)
		}
	})
}

func TestValidateStruct_CreateStudentRequest(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		req := model.CreateStudentRequest{
			NIM:         "187221000099",
			Nama:        "Test Student",
			Email:       "test@student.siakad.test",
			Prodi:       "Sistem Informasi",
			Angkatan:    2024,
			IPKTerakhir: 3.50,
		}
		if errs := ValidateStruct(req); errs != nil {
			t.Errorf("harap lolos, dapat: %v", errs)
		}
	})

	t.Run("NIM bukan 12 digit", func(t *testing.T) {
		req := model.CreateStudentRequest{
			NIM:      "123",
			Nama:     "Test Student",
			Email:    "test@student.siakad.test",
			Prodi:    "Sistem Informasi",
			Angkatan: 2024,
		}
		errs := ValidateStruct(req)
		if _, ok := errs["nim"]; !ok {
			t.Errorf("harap error pada nim, dapat: %v", errs)
		}
	})

	t.Run("IPK di luar rentang", func(t *testing.T) {
		req := model.CreateStudentRequest{
			NIM:         "187221000099",
			Nama:        "Test Student",
			Email:       "test@student.siakad.test",
			Prodi:       "Sistem Informasi",
			Angkatan:    2024,
			IPKTerakhir: 5.0,
		}
		errs := ValidateStruct(req)
		if _, ok := errs["ipk_terakhir"]; !ok {
			t.Errorf("harap error pada ipk_terakhir, dapat: %v", errs)
		}
	})
}

func TestValidateStruct_CreateEnrollmentRequest(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		req := model.CreateEnrollmentRequest{
			CourseID:      1,
			TahunAkademik: "2026/2027-Ganjil",
		}
		if errs := ValidateStruct(req); errs != nil {
			t.Errorf("harap lolos, dapat: %v", errs)
		}
	})

	t.Run("tahun akademik salah format", func(t *testing.T) {
		req := model.CreateEnrollmentRequest{
			CourseID:      1,
			TahunAkademik: "2026-2027",
		}
		errs := ValidateStruct(req)
		if _, ok := errs["tahun_akademik"]; !ok {
			t.Errorf("harap error pada tahun_akademik, dapat: %v", errs)
		}
	})

	t.Run("course_id nol", func(t *testing.T) {
		req := model.CreateEnrollmentRequest{
			CourseID:      0,
			TahunAkademik: "2026/2027-Ganjil",
		}
		errs := ValidateStruct(req)
		if _, ok := errs["course_id"]; !ok {
			t.Errorf("harap error pada course_id, dapat: %v", errs)
		}
	})
}