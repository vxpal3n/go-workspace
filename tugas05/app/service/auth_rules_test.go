package service

import (
	"testing"

	"tugas05/app/model"
)

func TestValidateRegister_Success(t *testing.T) {
	req := model.RegisterRequest{
		NIM:      "S001",
		Name:     "Thaariq",
		Email:    "thaariq@example.com",
		Grade:    85.5,
		Password: "rahasia123",
	}
	errs := ValidateRegister(req)
	if len(errs) != 0 {
		t.Errorf("seharusnya lolos, tapi dapat error: %v", errs)
	}
}

func TestValidateRegister_WeakPassword(t *testing.T) {
	cases := []struct {
		name     string
		password string
	}{
		{"terlalu pendek", "abc1"},
		{"tanpa angka", "rahasiaku"},
		{"tanpa huruf", "12345678"},
		{"terlalu umum", "password1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := model.RegisterRequest{
				NIM:      "S001",
				Name:     "Thaariq",
				Email:    "thaariq@example.com",
				Grade:    85.5,
				Password: tc.password,
			}
			errs := ValidateRegister(req)
			if _, ok := errs["password"]; !ok {
				t.Errorf("harap error pada field password, dapat: %v", errs)
			}
		})
	}
}

func TestValidateRegister_InvalidEmail(t *testing.T) {
	req := model.RegisterRequest{
		NIM:      "S001",
		Name:     "Thaariq",
		Email:    "bukan-email",
		Grade:    85.5,
		Password: "rahasia123",
	}
	errs := ValidateRegister(req)
	if _, ok := errs["email"]; !ok {
		t.Errorf("harap error pada field email, dapat: %v", errs)
	}
}

func TestValidateRegister_GradeOutOfRange(t *testing.T) {
	req := model.RegisterRequest{
		NIM:      "S001",
		Name:     "Thaariq",
		Email:    "thaariq@example.com",
		Grade:    150,
		Password: "rahasia123",
	}
	errs := ValidateRegister(req)
	if _, ok := errs["grade"]; !ok {
		t.Errorf("harap error pada field grade, dapat: %v", errs)
	}
}