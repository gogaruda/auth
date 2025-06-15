package service

import (
	"errors"
	"sql/internal/model"
	"sql/internal/repository"
)

type UserService interface {
	GetAll() ([]model.UserModel, error)
	GetByID(userID uint) (*model.UserModel, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) GetAll() ([]model.UserModel, error) {
	return s.repo.GetAll()
}

func (s *userService) GetByID(userID uint) (*model.UserModel, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, errors.New("Data tidak ditemukan")
	}

	return user, nil
}
