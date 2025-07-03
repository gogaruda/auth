package service

import (
	"context"
	"github.com/gogaruda/auth/internal/repository"
)

type UserService interface {
	MarkEmailVerified(ctx context.Context, userID string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) MarkEmailVerified(ctx context.Context, userID string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateIsVerified(ctx, user); err != nil {
		return err
	}

	return nil
}
