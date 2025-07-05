package service

import (
	"context"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/dto/request"
	"github.com/gogaruda/auth/internal/dto/response"
	"github.com/gogaruda/auth/internal/model"
	"github.com/gogaruda/auth/internal/repository"
	"github.com/gogaruda/auth/pkg/utils"
)

type UserService interface {
	Create(ctx context.Context, user *request.UserCreateRequest) error
	GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error)
	MarkEmailVerified(ctx context.Context, userID string) error
}

type userService struct {
	repo     repository.UserRepository
	authRepo repository.AuthRepository
	roleRepo repository.RoleRepository
	ut       utils.Utils
}

func NewUserService(r repository.UserRepository, auth repository.AuthRepository, role repository.RoleRepository, ut utils.Utils) UserService {
	return &userService{repo: r, authRepo: auth, roleRepo: role, ut: ut}
}

func (s *userService) Create(ctx context.Context, user *request.UserCreateRequest) error {
	usernameExists, err := s.authRepo.IsUsernameExists(ctx, user.Username)
	if err != nil {
		return err
	}
	if usernameExists {
		return apperror.New(apperror.CodeUsernameConflict, "username sudah digunakan", err)
	}

	emailExists, err := s.authRepo.IsEmailExists(ctx, user.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return apperror.New(apperror.CodeEmailConflict, "email sudah digunakan", err)
	}

	roles, err := s.roleRepo.CheckRoles(ctx, user.Roles)
	if err != nil {
		return err
	}

	hashPass, err := s.ut.GenerateHash(user.Password)
	if err != nil {
		return apperror.New(apperror.CodeInternalError, "gagal generate password", err)
	}

	userModel := model.UserModel{
		ID:             s.ut.GenerateULID(),
		Username:       &user.Username,
		Email:          user.Email,
		Password:       &hashPass,
		TokenVersion:   nil,
		GoogleID:       nil,
		IsVerified:     false,
		CreatedByAdmin: true,
		Roles:          roles,
	}

	if err := s.repo.Create(ctx, userModel); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error) {
	return s.repo.GetAll(ctx, limit, offset)
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
