package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/model"
	"github.com/gogaruda/auth/pkg/utils"
	"github.com/gogaruda/dbtx"
	"strings"
)

type AuthRepository interface {
	IsUsernameExists(ctx context.Context, username string) (bool, error)
	IsEmailExists(ctx context.Context, email string) (bool, error)
	CheckRoles(ctx context.Context, roles []string) ([]model.RoleModel, error)
	Create(ctx context.Context, user model.UserModel) error
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

func (r *authRepository) IsUsernameExists(ctx context.Context, username string) (bool, error) {
	const query = `SELECT exists (SELECT 1 FROM username_history WHERE username = ?)`

	var exists bool
	err := r.database.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, apperror.New(apperror.CodeDBError, "gagal memeriksa apakah username sudah terdaftar di database", err)
	}

	return exists, nil
}

func (r *authRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	const query = `SELECT exists (SELECT 1 FROM email_history WHERE email = ?)`

	var exists bool
	err := r.database.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, apperror.New(apperror.CodeDBError, "gagal memerika apakah email sudah terdaftar di datbase", err)
	}

	return exists, nil
}

func (r *authRepository) CheckRoles(ctx context.Context, roles []string) ([]model.RoleModel, error) {
	if len(roles) == 0 {
		return nil, apperror.New(apperror.CodeBadRequest, "roles tidak boleh kosong", errors.New("roles tidak bileh kosong"))
	}

	if len(roles) > 20 {
		return nil, apperror.New(apperror.CodeBadRequest, "jumlah roles terlalu banyak", errors.New("jumlah roles terlalu banyak"))
	}

	// Hilangkan duplikat role
	roleMap := make(map[string]struct{})
	var uniqueRoles []string
	for _, role := range roles {
		if _, exists := roleMap[role]; !exists {
			roleMap[role] = struct{}{}
			uniqueRoles = append(uniqueRoles, role)
		}
	}

	// Siapkan query
	placeholder := make([]string, len(uniqueRoles))
	args := make([]interface{}, len(uniqueRoles))
	for i, role := range uniqueRoles {
		placeholder[i] = "?"
		args[i] = role
	}

	query := fmt.Sprintf(
		`SELECT id, name FROM roles WHERE name IN (%s)`,
		strings.Join(placeholder, ","),
	)

	rows, err := r.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.New(apperror.CodeDBError, "query roles gagal", err)
	}
	defer rows.Close()

	var foundRoles []model.RoleModel
	for rows.Next() {
		var role model.RoleModel
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, apperror.New(apperror.CodeDBError, "gagal membaca data role", err)
		}
		foundRoles = append(foundRoles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.New(apperror.CodeDBError, "error saat membaca hasil query", err)
	}

	if len(foundRoles) != len(uniqueRoles) {
		return nil, apperror.New(apperror.CodeRoleNotFound, "salah satu atau lebih role tidak ditemukan", errors.New("salah satu atau lebih role tidak ditemukan"))
	}

	return foundRoles, nil
}

func (r *authRepository) Create(ctx context.Context, user model.UserModel) error {
	return dbtx.WithTxContext(ctx, r.database, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO users(id, username, email, password) VALUES(?, ?, ?, ?)`,
			user.ID, user.Username, user.Email, user.Password)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "query insert users gagal", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO username_history(username) VALUES(?)`, user.Username)
		if err != nil {
			return apperror.New(apperror.CodeDBError, "query insert username_history gagal", err)
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
