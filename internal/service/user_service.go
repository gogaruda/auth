package service

import (
	"sql/internal/dto/request"
	"sql/internal/dto/response"
	"sql/internal/repository"
)

type UserService interface {
	GetAll() ([]response.UserResponse, error)
	GetByID(userID string) (*response.UserResponse, error)
	UpdateUser(req request.UpdateUserRequest) error
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

func (s *userService) GetByID(userID string) (*response.UserResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(req request.UpdateUserRequest) error {
	return s.repo.Update(req)
}
