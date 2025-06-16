package repository

import (
	"database/sql"
	"sql/internal/model"
	"sql/pkg/utils"
)

type AuthRepository interface {
	IdentifierCheck(identifier string) (*model.UserModel, error)
	UpdateTokenVersion(userID string) (string, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) IdentifierCheck(identifier string) (*model.UserModel, error) {
	var user model.UserModel

	err := r.db.QueryRow(`
		SELECT id, username, email, password, token_version 
		FROM users 
		WHERE username = ? OR email = ?`,
		identifier, identifier,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.TokenVersion,
	)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT r.id, r.name
		FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`,
		user.ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.RoleModel
	for rows.Next() {
		var role model.RoleModel
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	user.Roles = roles

	return &user, nil
}

func (r *authRepository) UpdateTokenVersion(userID string) (string, error) {
	newVersion := utils.NewULID()
	_, err := r.db.Exec("UPDATE users SET token_version = ? WHERE id = ?", newVersion, userID)
	if err != nil {
		return "", err
	}

	return newVersion, nil
}
