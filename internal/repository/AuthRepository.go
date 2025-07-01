package repository

import (
	"context"
	"database/sql"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/model"
	"github.com/gogaruda/auth/pkg/utils"
	"github.com/gogaruda/dbtx"
)

type AuthRepository interface {
	Identifier(ctx context.Context, identifier string) (*model.UserModel, error)
	UpdateTokenVersion(userID string) (string, error)
}

type authRepository struct {
	database *sql.DB
	id       utils.ULIDs
}

func NewAuthRepository(db *sql.DB, i utils.ULIDs) AuthRepository {
	return &authRepository{database: db, id: i}
}

func (r *authRepository) Identifier(ctx context.Context, identifier string) (*model.UserModel, error) {
	var user model.UserModel
	var roles []model.RoleModel

	err := dbtx.WithTxContext(ctx, r.database, func(ctx context.Context, tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, `SELECT id, password FROM users WHERE username = ? OR email = ?`, identifier, identifier).
			Scan(&user.ID, &user.Password)
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

func (r *authRepository) UpdateTokenVersion(userID string) (string, error) {
	newVersion := r.id.Create()
	_, err := r.database.Exec(`UPDATE users SET token_version = ? WHERE id = ?`, newVersion, userID)
	if err != nil {
		return "", apperror.New(apperror.CodeDBError, "query update token_version gagal", err)
	}

	return newVersion, nil
}
