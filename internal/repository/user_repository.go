package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogaruda/auth/internal/dto/request"
	"github.com/gogaruda/auth/internal/dto/response"
	"github.com/gogaruda/auth/pkg/utils"
)

type UserRepository interface {
	IsRoleExists(roles []string) error
	GetAll() ([]response.UserResponse, error)
	Create(req request.CreateUserRequest) error
	GetByID(userID string) (*response.UserResponse, error)
	Update(userID string, req request.UpdateUserRequest) error
	Delete(userID string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) IsRoleExists(roles []string) error {
	for _, roleID := range roles {
		var exists bool
		err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM roles WHERE id = ?)`, roleID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("gagal query role_id %v: %w", roleID, err)
		}

		if !exists {
			return fmt.Errorf("role_id %v tidak ditemukan di database", roleID)
		}
	}

	return nil
}

func (r *userRepository) GetAll() ([]response.UserResponse, error) {
	query := `SELECT u.id, u.username, u.email, r.id AS role_id, r.name AS role
						FROM user_roles ur
						INNER JOIN users u ON ur.user_id = u.id
						INNER JOIN roles r ON ur.role_id = r.id
						WHERE r.name != ?
						ORDER BY u.updated_at DESC`
	rows, err := r.db.Query(query, "super-admin")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userMap := make(map[string]*response.UserResponse)

	for rows.Next() {
		var (
			userID, username, email string
			roleID, roleName        string
		)

		if err := rows.Scan(&userID, &username, &email, &roleID, &roleName); err != nil {
			return nil, err
		}

		if _, exists := userMap[userID]; !exists {
			userMap[userID] = &response.UserResponse{
				ID:       userID,
				Username: username,
				Email:    email,
				Roles:    []response.RoleResponse{},
			}
		}

		userMap[userID].Roles = append(userMap[userID].Roles, response.RoleResponse{
			ID:   roleID,
			Name: roleName,
		})
	}

	var users []response.UserResponse
	for _, user := range userMap {
		users = append(users, *user)
	}

	return users, nil
}

func (r *userRepository) Create(req request.CreateUserRequest) error {
	newUserID := utils.NewULID()
	hashedPassword, err := utils.GenerateHash(req.Password)
	if err != nil {
		return errors.New("gagal generate password")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.Exec(`INSERT INTO users(id, username, email, password) VALUES(?, ?, ?, ?)`,
		newUserID,
		req.Username,
		req.Email,
		hashedPassword,
	)
	if err != nil {
		return errors.New("gagal tambah data users")
	}

	stmt, err := tx.Prepare(`INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, roleID := range req.RoleIDs {
		if _, err := stmt.Exec(newUserID, roleID); err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) GetByID(userID string) (*response.UserResponse, error) {
	query := `
		SELECT 
			u.id, u.username, u.email,
			r.id AS role_id, r.name AS role_name
		FROM 
			user_roles ur
		JOIN users u ON ur.user_id = u.id
		JOIN roles r ON ur.role_id = r.id
		WHERE u.id = ?
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user *response.UserResponse

	for rows.Next() {
		var (
			id, username, email string
			roleID, roleName    string
		)

		if err := rows.Scan(&id, &username, &email, &roleID, &roleName); err != nil {
			return nil, err
		}

		if user == nil {
			user = &response.UserResponse{
				ID:       id,
				Username: username,
				Email:    email,
				Roles:    []response.RoleResponse{},
			}
		}

		user.Roles = append(user.Roles, response.RoleResponse{
			ID:   roleID,
			Name: roleName,
		})
	}

	if user == nil {
		return nil, errors.New("User tidak ditemukan!")
	}

	return user, nil
}

func (r *userRepository) Update(userID string, req request.UpdateUserRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.Exec(`UPDATE users SET username = ?, email = ? WHERE id = ?`,
		req.Username, req.Email, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM user_roles WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, roleID := range req.RoleIDs {
		if _, err := stmt.Exec(userID, roleID); err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) Delete(userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.Exec(`DELETE FROM users WHERE id = ?`, userID)
	if err != nil {
		return errors.New("Gagal delete user")
	}

	_, err = tx.Exec(`DELETE FROM user_roles WHERE user_id = ?`, userID)
	if err != nil {
		return errors.New("Gagal delete user_roles")
	}

	return nil
}
