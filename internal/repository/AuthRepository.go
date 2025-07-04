package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/model"
	"github.com/gogaruda/dbtx"
)

type AuthRepository interface {
	IsUsernameExists(ctx context.Context, username string) (bool, error)
	IsEmailExists(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, user model.UserModel) error
	Identifier(ctx context.Context, identifier string) (*model.UserModel, error)
	UpdateTokenVersion(userID, newVersion string) error
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	const query = `SELECT exists (SELECT 1 FROM username_history WHERE username = ?)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, apperror.New(apperror.CodeDBError, "gagal memeriksa apakah username sudah terdaftar di db", err)
	}

	return exists, nil
}

func (r *authRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	const query = `SELECT exists (SELECT 1 FROM email_history WHERE email = ?)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, apperror.New(apperror.CodeDBError, "gagal memerika apakah email sudah terdaftar di datbase", err)
	}

	return exists, nil
}

func (r *authRepository) Create(ctx context.Context, user model.UserModel) error {
	return dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO 
			users(id, username, email, password, token_version, google_id, is_verified, created_by_admin) 
			VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
			user.ID, user.Username, user.Email, user.Password,
			user.TokenVersion, user.GoogleID, user.IsVerified, user.CreatedByAdmin)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "query insert users gagal", err)
		}

		if user.Username != nil {
			_, err = tx.ExecContext(ctx, `INSERT INTO username_history(username) VALUES(?)`, user.Username)
			if err != nil {
				return apperror.New(apperror.CodeDBError, "query insert username_history gagal", err)
			}
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO email_history(email) VALUES(?)`, user.Email)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "query insert email_history gagal", err)
		}

		stmt, err := tx.PrepareContext(ctx, `INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)`)
		if err != nil {
			return apperror.New(apperror.CodeDBPrepareError, "gagal prepare insert user_roles", err)
		}
		defer stmt.Close()

		for _, r := range user.Roles {
			_, err := stmt.ExecContext(ctx, user.ID, r.ID)
			if err != nil {
				return apperror.New(apperror.CodeDBError,
					fmt.Sprintf("query insert user_roles gagal untuk role_id: %s", r.ID), err)
			}
		}

		return nil
	})
}

func (r *authRepository) Identifier(ctx context.Context, identifier string) (*model.UserModel, error) {
	var user model.UserModel
	var roles []model.RoleModel

	err := dbtx.WithTxContext(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, `SELECT id, password, is_verified FROM users WHERE username = ? OR email = ?`, identifier, identifier).
			Scan(&user.ID, &user.Password, &user.IsVerified)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "query select users gagal", err)
		}

		query := `SELECT r.id, r.name FROM roles r INNER JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = ?`
		rows, err := tx.QueryContext(ctx, query, user.ID)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "query select roles gagal", err)
		}
		defer rows.Close()

		for rows.Next() {
			var role model.RoleModel
			if err := rows.Scan(&role.ID, &role.Name); err != nil {
				return apperror.New(apperror.CodeDBError, "gagal scan roles", err)
			}
			roles = append(roles, role)
		}

		if err := rows.Err(); err != nil {
			return apperror.New(apperror.CodeDBError, "gagal setelah iterasi", err)
		}

		user.Roles = roles
		return nil
	})

	if err != nil {
		return nil, apperror.New(apperror.CodeDBTxFailed, "gagal query context", err)
	}

	return &user, nil
}

func (r *authRepository) UpdateTokenVersion(userID, newVersion string) error {
	_, err := r.db.Exec(`UPDATE users SET token_version = ? WHERE id = ?`, newVersion, userID)
	if err != nil {
		return apperror.New(apperror.CodeDBError, "query update token_version gagal", err)
	}

	return nil
}
