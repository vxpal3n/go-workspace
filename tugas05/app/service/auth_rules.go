package service

import (
	"strings"
	"unicode"

	"tugas05/app/model"
)

const minPasswordLength = 8

func ValidateRegister(req model.RegisterRequest) map[string]string {
	errs := map[string]string{}

	nim := strings.TrimSpace(req.NIM)
	switch {
	case nim == "":
		errs["nim"] = "wajib diisi"
	case len(nim) < 3:
		errs["nim"] = "minimal 3 karakter"
	}

	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "wajib diisi"
	}

	if !isValidEmail(req.Email) {
		errs["email"] = "format email tidak valid"
	}

	if req.Grade < 0 || req.Grade > 100 {
		errs["grade"] = "harus antara 0 dan 100"
	}

	if msg := checkPasswordStrength(req.Password); msg != "" {
		errs["password"] = msg
	}
	return errs
}

func ValidateLogin(req model.LoginRequest) map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "wajib diisi"
	}
	if req.Password == "" {
		errs["password"] = "wajib diisi"
	}
	return errs
}

func checkPasswordStrength(password string) string {
	if len(password) < minPasswordLength {
		return "minimal 8 karakter"
	}
	var hasLetter, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return "harus memuat huruf dan angka"
	}
	weak := map[string]bool{
		"password1": true, "12345678": true, "qwerty123": true,
		"admin123": true, "password123": true,
	}
	if weak[strings.ToLower(password)] {
		return "password terlalu umum"
	}
	return ""
}

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at+1 && dot < len(email)-1
}