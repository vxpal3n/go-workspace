package service

import "strings"

func ValidateLoginEmail(email string) string {
	if strings.TrimSpace(email) == "" {
		return "wajib diisi"
	}
	return ""
}