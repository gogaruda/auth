package repository

import (
	"database/sql"
	"sql/internal/model"
)

type UserRepository interface {
	GetAll() ([]model.UserModel, error)
	GetByID(userID uint) (*model.UserModel, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll() ([]model.UserModel, error) {
	var users []model.UserModel
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user model.UserModel
		if err := rows.Scan(&user.ID, &user.Nama, &user.Alamat); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, err
}

func (r *userRepository) GetByID(userID uint) (*model.UserModel, error) {
	var user model.UserModel
	err := r.db.QueryRow("SELECT id, nama, alamat FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Nama, &user.Alamat)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
