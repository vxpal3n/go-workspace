package service

import (
	"strings"

	"tugas07/app/model"
	"tugas07/helper"
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

	if msg := helper.CheckPasswordStrength(req.Password); msg != "" {
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

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at+1 && dot < len(email)-1
}
