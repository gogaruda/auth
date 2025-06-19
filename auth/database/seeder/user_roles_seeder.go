package seeder

import (
	"database/sql"
	"github.com/gogaruda/auth/pkg/utils"
)

func User_roles(db *sql.DB) error {
	passwordHash, err := utils.GenerateHash("super-admin-1")
	if err != nil {
		return err
	}

	userID := utils.NewULID()
	superAdminRoleID := utils.NewULID()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO users (id, username, email, password) 
		VALUES (?, ?, ?, ?)`,
		userID, "super-admin-1", "super-admin-1@gmail.com", passwordHash,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO roles (id, name) VALUES 
		(?, ?), (?, ?), (?, ?), (?, ?), (?, ?)`,
		superAdminRoleID, "super-admin",
		utils.NewULID(), "admin",
		utils.NewULID(), "editor",
		utils.NewULID(), "penulis",
		utils.NewULID(), "tamu",
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`,
		userID, superAdminRoleID,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
