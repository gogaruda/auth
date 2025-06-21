package service

import (
	"github.com/gogaruda/auth/auth/dto/request"
	"github.com/gogaruda/auth/auth/repository"
	"github.com/gogaruda/pkg/apperror"
	"github.com/gogaruda/pkg/utils"
)

type AuthService interface {
	Login(request request.AuthLoginRequest) (string, error)
	Register(req request.AuthRegisterRequest) error
	Logout(userID string) error
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return &authService{repo: r}
}

func (s *authService) Login(request request.AuthLoginRequest) (string, error) {
	user, err := s.repo.IdentifierCheck(request.Identifier)
	if err != nil || !utils.CompareHash(user.Password, request.Password) {
		return "", apperror.New(apperror.CodeAuthNotFound, "username/email atau password salah", err)
	}

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}

	newTokenVersion, err := s.repo.UpdateTokenVersion(user.ID)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateJWT(user.ID, newTokenVersion, roleNames)

	return token, err
}

func (s *authService) Register(req request.AuthRegisterRequest) error {
	// Cek username
	exists, err := s.repo.IsUsernameExists(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return apperror.New(apperror.CodeUserConflict, "username sudah terdaftar", err)
	}

	// Cek email
	existsEmail, errEmail := s.repo.IsEmailExists(req.Email)
	if errEmail != nil {
		return err
	}
	if existsEmail {
		return apperror.New(apperror.CodeUserConflict, "email sudah terdaftar", err)
	}

	// Create user
	if err := s.repo.Create(req); err != nil {
		return err
	}

	return nil
}

func (s *authService) Logout(userID string) error {
	_, err := s.repo.UpdateTokenVersion(userID)
	if err != nil {
		return err
	}
	return nil
}
