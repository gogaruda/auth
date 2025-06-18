package service

import (
	"github.com/gogaruda/auth/internal/dto/request"
	"github.com/gogaruda/auth/internal/dto/response"
	"github.com/gogaruda/auth/internal/repository"
)

type UserService interface {
	GetAll() ([]response.UserResponse, error)
	Create(req request.CreateUserRequest) error
	GetByID(userID string) (*response.UserResponse, error)
	UpdateUser(userID string, req request.UpdateUserRequest) error
	Delete(userID string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetAll() ([]response.UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) Create(req request.CreateUserRequest) error {
	if err := s.repo.IsRoleExists(req.RoleIDs); err != nil {
		return err
	}

	if err := s.repo.Create(req); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetByID(userID string) (*response.UserResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(userID string, req request.UpdateUserRequest) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	return s.repo.Update(user.ID, req)
}

func (s *userService) Delete(userID string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(user.ID); err != nil {
		return err
	}

	return nil
}
