package seeder

import (
	"database/sql"
	"sql/pkg/utils"
)

func User_roles(db *sql.DB) error {
	passwordHash, err := utils.GenerateHash("super-admin-1")
	if err != nil {
		return err
	}

	userID := utils.NewULID()
	roleID := utils.NewULID()

	_, err = db.Exec("INSERT INTO users (id, username, email, password) VALUES(?, ?, ?, ?)",
		userID, "super-admin-1", "super-admin-1@gmail.com", passwordHash,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO roles (id, name) VALUES(?, ?)",
		roleID, "super-admin",
	)

	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)",
		userID, roleID,
	)

	if err != nil {
		return err
	}

	return nil
}
