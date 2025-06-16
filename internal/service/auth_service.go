package service

import (
	"errors"
	"sql/internal/dto/request"
	"sql/internal/repository"
	"sql/pkg/utils"
)

type AuthService interface {
	Login(request request.AuthLoginRequest) (string, error)
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
		return "", errors.New("username/email atau password salah")
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
