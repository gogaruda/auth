package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/auth/internal/model"
	"strings"
)

type RoleRepository interface {
	CheckRoles(ctx context.Context, roles []string) ([]model.RoleModel, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) CheckRoles(ctx context.Context, roles []string) ([]model.RoleModel, error) {
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

	rows, err := r.db.QueryContext(ctx, query, args...)
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
