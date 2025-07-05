package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/dto/response"
	"github.com/gogaruda/auth/internal/model"
	"github.com/gogaruda/dbtx"
)

type UserRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error)
	Create(ctx context.Context, user model.UserModel) error
	FindByEmail(ctx context.Context, email string) (*model.UserModel, error)
	FindByID(ctx context.Context, userID string) (*model.UserModel, error)
	UpdateIsVerified(ctx context.Context, user *model.UserModel) error
	UpdateGoogleID(ctx context.Context, userID, googleID string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]response.UserResponse, int, error) {
	var total int
	queryCount := `
				SELECT COUNT(DISTINCT u.id) AS total
				FROM users u
				JOIN user_roles ur ON u.id = ur.user_id
				JOIN roles r ON r.id = ur.role_id
				WHERE r.name NOT IN (?, ?)
			`
	err := r.db.QueryRowContext(ctx, queryCount, "super admin", "admin").Scan(&total)
	if err != nil {
		return nil, 0, apperror.New(apperror.CodeDBError, "query count users gagal", err)
	}

	queryUsers := `
		SELECT 
			u.id, u.username, u.email, u.google_id, u.is_verified, u.created_by_admin,
			r.id AS role_id, r.name AS role_name
		FROM users u 
		JOIN user_roles ur ON u.id = ur.user_id
		JOIN roles r ON r.id = ur.role_id
		WHERE u.id NOT IN (
			SELECT ur.user_id 
			FROM user_roles ur
			JOIN roles r ON r.id = ur.role_id
			WHERE r.name IN (?,?)
		)
		ORDER BY u.updated_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, queryUsers, "super admin", "admin", limit, offset)
	if err != nil {
		return nil, 0, apperror.New(apperror.CodeDBError, "query select users gagal", err)
	}
	defer rows.Close()

	userMap := make(map[string]*response.UserResponse)
	var orderedIDs []string

	for rows.Next() {
		var (
			id, roleID, roleName       string
			username, googleID         sql.NullString
			email                      string
			isVerified, createdByAdmin bool
		)

		if err := rows.Scan(&id, &username, &email, &googleID, &isVerified, &createdByAdmin, &roleID, &roleName); err != nil {
			return nil, 0, apperror.New(apperror.CodeDBError, "gagal scan users", err)
		}

		user, exists := userMap[id]
		if !exists {
			user = &response.UserResponse{
				ID:             id,
				Email:          email,
				IsVerified:     isVerified,
				CreatedByAdmin: createdByAdmin,
				Roles:          []response.RoleResponse{},
			}

			if username.Valid {
				user.Username = &username.String
			}
			if googleID.Valid {
				user.GoogleID = &googleID.String
			}

			userMap[id] = user
			orderedIDs = append(orderedIDs, id)
		}

		user.Roles = append(user.Roles, response.RoleResponse{
			ID:   roleID,
			Name: roleName,
		})
	}

	var users []response.UserResponse
	for _, id := range orderedIDs {
		users = append(users, *userMap[id])
	}

	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user model.UserModel) error {
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

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	query := `SELECT id, username, email, password, google_id, is_verified, created_by_admin FROM users WHERE email = ?`
	var user model.UserModel
	var username, password, googleID sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &username, &user.Email, &password, &googleID, &user.IsVerified, &user.CreatedByAdmin,
	)

	if err != nil {
		return nil, err
	}

	if username.Valid {
		user.Username = &username.String
	}

	if password.Valid {
		user.Password = &password.String
	}

	if googleID.Valid {
		user.GoogleID = &googleID.String
	}

	rolesQuery := `SELECT r.id, r.name FROM roles r INNER JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = ?`
	rows, err := r.db.QueryContext(ctx, rolesQuery, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role model.RoleModel
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)
	}

	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, userID string) (*model.UserModel, error) {
	var user model.UserModel
	var tokenVersion sql.NullString
	user.TokenVersion = &tokenVersion.String
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

func (r *userRepository) UpdateGoogleID(ctx context.Context, googleID, userID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET google_id = ? WHERE id = ?`, googleID, userID)
	if err != nil {
		return apperror.New(apperror.CodeDBError, "googleID gagal di update", err)
	}
	return nil
}
