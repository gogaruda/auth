package repository

import (
	"database/sql"
	"sql/internal/dto/request"
	"sql/internal/model"
	"sql/pkg/utils"
)

type AuthRepository interface {
	IdentifierCheck(identifier string) (*model.UserModel, error)
	UpdateTokenVersion(userID string) (string, error)
	IsUsernameExists(username string) (bool, error)
	IsEmailExists(email string) (bool, error)
	Create(req request.AuthRegisterRequest) error
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

func (r *authRepository) IsUsernameExists(username string) (bool, error) {
	var existingUsername string
	err := r.db.QueryRow("SELECT username FROM users WHERE username = ?", username).
		Scan(&existingUsername)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (r *authRepository) IsEmailExists(email string) (bool, error) {
	var existingEmail string
	err := r.db.QueryRow("SELECT email FROM users WHERE email = ?", email).
		Scan(&existingEmail)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (r *authRepository) Create(req request.AuthRegisterRequest) error {
	userID := utils.NewULID()
	hashedPassword, err := utils.GenerateHash(req.Password)
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert user
	_, err = tx.Exec(`
			INSERT INTO 
			users(id, username, email, password) 
			VALUES (?, ?, ?, ?)
			`,
		userID,
		req.Username,
		req.Email,
		hashedPassword,
	)
	if err != nil {
		return err
	}

	// Ambil role_id
	var roleID string
	err = tx.QueryRow(`SELECT id FROM roles WHERE name = ?`, "tamu").Scan(&roleID)
	if err == sql.ErrNoRows {
		return err
	}

	// Insert relasi
	_, err = tx.Exec(`
			INSERT INTO
			user_roles(user_id, role_id)
			VALUES(?, ?)
			`,
		userID,
		roleID,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
