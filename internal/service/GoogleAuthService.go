package service

import (
	"context"
	"database/sql"
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
	cfg      *config.AppConfig
	ut       utils.Utils
}

func NewGoogleAuthService(r repository.UserRepository, c *config.AppConfig, u utils.Utils) GoogleAuthService {
	return &googleAuthService{userRepo: r, cfg: c, ut: u}
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

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Buat user baru
		if err == sql.ErrNoRows {
			user = &model.UserModel{
				ID:             s.ut.GenerateULID(),
				Username:       s.ut.GenerateUsernameFromName(userInfo.Name),
				Email:          userInfo.Email,
				Password:       nil,
				TokenVersion:   s.ut.GenerateULID(),
				GoogleID:       &userInfo.Id,
				IsVerified:     true,
				CreatedByAdmin: false,
			}

			if err := s.userRepo.Create(ctx, *user); err != nil {
				return "", err
			}
		}
		return "", err
	} else {
		// User sudah ada
		if user.GoogleID == nil {
			user.GoogleID = &userInfo.Id
			if err := s.userRepo.UpdateGoogleID(ctx, user.ID, *user.GoogleID); err != nil {
				return "", err
			}
		}

		// Jika user belum verifikasi email tolak login
		if user.CreatedByAdmin && !user.IsVerified {
			return "", apperror.New(apperror.CodeForbidden, "akun harus verifikasi email terlebih dahulu", err)
		}
	}

	return "", nil
}
