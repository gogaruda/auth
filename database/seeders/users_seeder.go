package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/dbtx"
	"github.com/gogaruda/pkg/utils"
	"time"
)

func Users(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return dbtx.WithTxContext(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		var roleID string
		err := tx.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = ?`, "super admin").Scan(&roleID)
		if err != nil {
			return fmt.Errorf("query roles gagal: %w", err)
		}

		userID := utils.NewULID()
		passHash, _ := utils.GenerateHash("superadmin")
		_, err = tx.ExecContext(ctx, `INSERT INTO users(id, username, email, password) VALUES(?, ?, ?, ?)`,
			userID, "superadmin", "superadmin@gmail.com", passHash)
		if err != nil {
			return fmt.Errorf("query insert users gagal: %w", err)
		}

		stmt, err := tx.PrepareContext(ctx, `INSERT INTO user_roles(user_id, role_id) VALUES(?, ?)`)
		if err != nil {
			return fmt.Errorf("prepare query user_roles gagal: %w", err)
		}
		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, userID, roleID)
		if err != nil {
			return fmt.Errorf("query insert user_roles gagal: %w", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO username_history(username) VALUES(?)`, "superadmin")
		if err != nil {
			return fmt.Errorf("query insert username history gagal: %w", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO email_history(email) VALUES(?)`, "superadmin@gmail.com")
		if err != nil {
			return fmt.Errorf("query insert email history gagal: %w", err)
		}

		_, err = tx.ExecContext(ctx, `INSERT INTO profiles(id, user_id, full_name, address, gender) VALUES(?, ?, ?, ?, ?)`,
			utils.NewULID(), userID, "Saya Super Admin Pertama", "Samarang - Garut", 1)
		if err != nil {
			return fmt.Errorf("query insert profiles gagal: %w", err)
		}
    
		return nil
	})
}
