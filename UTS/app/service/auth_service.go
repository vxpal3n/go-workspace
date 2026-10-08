package service

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"UTS/app/model"
	"UTS/app/repository"
	"UTS/helper"
)

type AuthService struct {
	users    repository.UserRepository
	students repository.StudentRepository
	jwt      *helper.JWTManager
}

func NewAuthService(
	users repository.UserRepository,
	students repository.StudentRepository,
	jwtManager *helper.JWTManager,
) *AuthService {
	return &AuthService{users: users, students: students, jwt: jwtManager}
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body harus berupa JSON yang valid")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	user, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		helper.VerifyDummyPassword(req.Password)
		return helper.Unauthorized("email atau password salah")
	}
	if !helper.VerifyPassword(user.Password, req.Password) {
		return helper.Unauthorized("email atau password salah")
	}

	accessToken, err := s.jwt.GenerateAccess(user)
	if err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, fiber.StatusOK, "Login berhasil", model.LoginData{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.jwt.AccessTTL().Seconds()),
		User: model.UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func (s *AuthService) Me(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("belum terautentikasi")
	}

	user, err := s.users.FindByID(ctx, authUser.UserID)
	if err != nil {
		return helper.Unauthorized("user tidak ditemukan")
	}

	resp := model.MeResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	if user.Role == "mahasiswa" {
		student, err := s.students.FindByUserID(ctx, user.ID)
		if err == nil {
			resp.Student = &model.MeStudentInfo{
				NIM:      student.NIM,
				Nama:     student.Nama,
				Prodi:    student.Prodi,
				Angkatan: student.Angkatan,
			}
		}
	}

	return helper.Success(c, fiber.StatusOK, "Profil berhasil diambil", resp)
}