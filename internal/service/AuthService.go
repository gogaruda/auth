package service

import (
	"context"
	"errors"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/internal/dto/request"
	"github.com/gogaruda/auth/internal/model"
	"github.com/gogaruda/auth/internal/repository"
	"github.com/gogaruda/auth/pkg/utils"
)

type AuthService interface {
	Register(ctx context.Context, req request.RegisterRequest) error
	Login(ctx context.Context, req request.LoginRequest) (string, error)
	Logout(userID string) error
}

type authService struct {
	authRepo repository.AuthRepository
	config   *config.AppConfig
	hash     utils.Hash
	jwt      utils.JWTs
	id       utils.ULIDs
	email    EmailVerificationService
}

func NewAuthService(
	a repository.AuthRepository,
	cfg *config.AppConfig,
	h utils.Hash,
	j utils.JWTs,
	i utils.ULIDs,
	e EmailVerificationService,
) AuthService {
	return &authService{authRepo: a, config: cfg, hash: h, jwt: j, id: i, email: e}
}

func (s *authService) Register(ctx context.Context, req request.RegisterRequest) error {
	isUsernameExists, err := s.authRepo.IsUsernameExists(ctx, req.Username)
	if err != nil {
		return err
	}
	if isUsernameExists {
		return apperror.New(apperror.CodeUsernameConflict, "username sudah terdaftar", errors.New("username sudah terdaftar"))
	}

	isEmailExists, err := s.authRepo.IsEmailExists(ctx, req.Email)
	if err != nil {
		return err
	}
	if isEmailExists {
		return apperror.New(apperror.CodeEmailConflict, "email sudah terdaftar", errors.New("email sudah terdaftar"))
	}

	roles, err := s.authRepo.CheckRoles(ctx, req.Roles)
	if err != nil {
		return err
	}

	hashPass, _ := s.hash.Generate(req.Password)
	user := model.UserModel{
		ID:       s.id.Create(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashPass,
		Roles:    roles,
	}

	if err := s.authRepo.Create(ctx, user); err != nil {
		return err
	}

	if err := s.email.SendVerification(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(ctx context.Context, req request.LoginRequest) (string, error) {
	user, err := s.authRepo.Identifier(ctx, req.Identifier)
	if err != nil || !s.hash.Compare(user.Password, req.Password) {
		return "", err
	}

	newVersion := s.id.Create()
	if err := s.authRepo.UpdateTokenVersion(user.ID, newVersion); err != nil {
		return "", err
	}

	var roles []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	token, err := s.jwt.Create(user.ID, newVersion, roles, s.config)
	if err != nil {
		return "", apperror.New(apperror.CodeInternalError, "gagal buat JWT", err)
	}

	return token, nil
}

func (s *authService) Logout(userID string) error {
	if err := s.authRepo.UpdateTokenVersion(userID, s.id.Create()); err != nil {
		return err
	}

	return nil
}
