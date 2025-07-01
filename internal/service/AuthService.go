package service

import (
	"context"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/internal/dto/request"
	"github.com/gogaruda/auth/internal/repository"
	"github.com/gogaruda/auth/pkg/utils"
)

type AuthService interface {
	Login(ctx context.Context, req request.LoginRequest) (string, error)
}

type authService struct {
	authRepo repository.AuthRepository
	config   *config.AppConfig
	hash     utils.Hash
	jwt      utils.JWTs
}

func NewAuthService(a repository.AuthRepository, cfg *config.AppConfig, h utils.Hash, j utils.JWTs) AuthService {
	return &authService{authRepo: a, config: cfg, hash: h, jwt: j}
}

func (s *authService) Login(ctx context.Context, req request.LoginRequest) (string, error) {
	user, err := s.authRepo.Identifier(ctx, req.Identifier)
	if err != nil || !s.hash.Compare(user.Password, req.Password) {
		return "", err
	}

	newVersion, err := s.authRepo.UpdateTokenVersion(user.ID)
	if err != nil {
		return "", err
	}

	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.ID)
	}

	token, err := s.jwt.Create(user.ID, newVersion, roles, s.config)
	if err != nil {
		return "", apperror.New(apperror.CodeInternalError, "gagal buat JWT", err)
	}

	return token, nil
}
