package repository

import (
	"database/sql"
	"sql/internal/dto/request"
	"sql/internal/dto/response"
)

type UserRepository interface {
	GetAll() ([]response.UserResponse, error)
	GetByID(userID string) (*response.UserResponse, error)
	Update(userID string, req request.UpdateUserRequest) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
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
		return nil, sql.ErrNoRows
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
