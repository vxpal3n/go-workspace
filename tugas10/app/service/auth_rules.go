package service

import (
	"strings"

	"tugas10/app/model"
)

const minPasswordLength = 8

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
