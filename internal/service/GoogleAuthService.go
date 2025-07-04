package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/internal/model"
	"github.com/gogaruda/auth/internal/repository"
	"github.com/gogaruda/auth/pkg/utils"
	"google.golang.org/api/oauth2/v1"
)

type GoogleAuthService interface {
	Login() string
	Callback(ctx context.Context, code string) (string, error)
}

type googleAuthService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	authRepo repository.AuthRepository
	cfg      *config.AppConfig
	ut       utils.Utils
}

func NewGoogleAuthService(
	r repository.UserRepository,
	rr repository.RoleRepository,
	au repository.AuthRepository,
	c *config.AppConfig,
	u utils.Utils) GoogleAuthService {
	return &googleAuthService{
		userRepo: r,
		roleRepo: rr,
		cfg:      c,
		ut:       u, authRepo: au}
}

func (s *googleAuthService) Login() string {
	return s.cfg.Google.AuthCodeURL("state-random")
}

func (s *googleAuthService) Callback(ctx context.Context, code string) (string, error) {
	token, err := s.cfg.Google.Exchange(ctx, code)
	if err != nil {
		return "", apperror.New(apperror.CodeInternalError, "exchange token tidak valid", err)
	}

	client := s.cfg.Google.Client(ctx, token)
	service, err := oauth2.New(client)
	if err != nil {
		return "", apperror.New(apperror.CodeInternalError, "gagal membuat OAuth2 service", err)
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return "", apperror.New(apperror.CodeInternalError, "gagal mendapatkan user info", err)
	}

	email := userInfo.Email
	var finalUser *model.UserModel
	newTokenVersion := s.ut.GenerateULID()

	// Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(ctx, email)
	switch {
	case err == nil:
		// User sudah terdaftar
		if user.GoogleID == nil {
			user.GoogleID = &userInfo.Id
			if err := s.userRepo.UpdateGoogleID(ctx, user.ID, *user.GoogleID); err != nil {
				return "", err
			}
		}

		if user.CreatedByAdmin && !user.IsVerified {
			return "", apperror.New("[EMAIL_NOT_VERIFIED]", "akun harus verifikasi email terlebih dahulu", nil, 403)
		}

		if err := s.authRepo.UpdateTokenVersion(user.ID, newTokenVersion); err != nil {
			return "", err
		}

		finalUser = user

	case errors.Is(err, sql.ErrNoRows):
		tamuRole, err := s.roleRepo.CheckRoles(ctx, []string{"tamu"})
		if err != nil {
			return "", err
		}

		newUser := model.UserModel{
			ID:             s.ut.GenerateULID(),
			Username:       nil,
			Email:          email,
			Password:       nil,
			TokenVersion:   &newTokenVersion,
			GoogleID:       &userInfo.Id,
			IsVerified:     true,
			CreatedByAdmin: false,
			Roles:          tamuRole,
		}

		if err := s.userRepo.Create(ctx, newUser); err != nil {
			return "", err
		}

		finalUser = &newUser

	default:
		return "", apperror.New(apperror.CodeDBError, "gagal mencari user", err)
	}

	// ambil roles
	var roles []string
	for _, r := range finalUser.Roles {
		roles = append(roles, r.Name)
	}

	finalUser.TokenVersion = &newTokenVersion

	tokenString, err := s.ut.GenerateJWT(finalUser.ID, newTokenVersion, finalUser.IsVerified, roles, s.cfg)
	if err != nil {
		return "", apperror.New(apperror.CodeInternalError, "gagal buat JWT", err)
	}

	return tokenString, nil
}
