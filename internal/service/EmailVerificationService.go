package service

import (
	"context"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/internal/model"
	"github.com/gogaruda/auth/internal/repository"
	"github.com/gogaruda/auth/pkg/mailer"
	"github.com/gogaruda/auth/pkg/utils"
	"time"
)

type EmailVerificationService interface {
	SendVerification(ctx context.Context, user model.UserModel) error
	VerifyToken(ctx context.Context, token string) error
}

type emailVerificationService struct {
	evRepo      repository.EmailVerificationRepository
	mail        mailer.Mailer
	id          utils.ULIDs
	cMail       config.EmailConfig
	userService UserService
}

func NewEmailVerificationService(
	ev repository.EmailVerificationRepository,
	t mailer.Mailer,
	i utils.ULIDs,
	c config.EmailConfig,
	u UserService,
) EmailVerificationService {
	return &emailVerificationService{
		evRepo:      ev,
		mail:        t,
		id:          i,
		cMail:       c,
		userService: u,
	}
}

func (s *emailVerificationService) SendVerification(ctx context.Context, user model.UserModel) error {
	tok, err := s.mail.GenerateRandomToken(32)
	if err != nil {
		return err
	}

	ev := &model.EmailVerificationModel{
		ID:        s.id.Create(),
		UserID:    user.ID,
		Token:     tok,
		ExpiresAt: time.Now().UTC().Add(30 * time.Minute),
	}

	if err := s.evRepo.Create(ctx, ev); err != nil {
		return err
	}

	url := fmt.Sprintf("%s?token=%s", s.cMail.FrontVerifyUrl, tok)
	body := fmt.Sprintf("<p>Klik disini untuk verifikasi email: <a href='%s'>%s</a></p>", url, url)

	if err := s.mail.Send(user.Email, "Verifikasi Email", body); err != nil {
		return apperror.New(apperror.CodeInternalError, "gagal mengirim verifikasi email", err)
	}
	return nil
}

func (s *emailVerificationService) VerifyToken(ctx context.Context, token string) error {
	ev, err := s.evRepo.FindByToken(ctx, token)
	if err != nil {
		return err
	}

	if ev.IsUsed {
		return apperror.New("[TOKEN_USED]", "token sudah digunakan", err, 505)
	}

	if time.Now().UTC().After(ev.ExpiresAt) {
		return apperror.New("[TOKEN_EXPIRED]", "token sudah kadaluarsa", err, 505)
	}

	if err := s.userService.MarkEmailVerified(ctx, ev.UserID); err != nil {
		return err
	}

	return s.evRepo.MarkAsUsed(ctx, ev.ID)
}
