package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/model"
)

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (*model.UserModel, error)
	UpdateIsVerified(ctx context.Context, user *model.UserModel) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, userID string) (*model.UserModel, error) {
	var user model.UserModel
	var tokenVersion sql.NullString
	user.TokenVersion = tokenVersion.String
	err := r.db.QueryRowContext(ctx, `SELECT id, username, email, password, token_version, is_verified FROM users WHERE id = ? LIMIT 1`, userID).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password, &tokenVersion, &user.IsVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.New(apperror.CodeUserNotFound, "user tidak ditemukan", err)
		}
		return nil, apperror.New(apperror.CodeDBError, "query findbyid users gagal", err)
	}

	return &user, nil
}

func (r *userRepository) UpdateIsVerified(ctx context.Context, user *model.UserModel) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET is_verified = true WHERE id = ?`, user.ID)
	if err != nil {
		return apperror.New(apperror.CodeDBError, "query update users is_verified gagal", err)
	}

	return nil
}
