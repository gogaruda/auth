package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogaruda/auth/pkg/utils"
	"github.com/gogaruda/dbtx"
	"time"
)

func Roles(db *sql.DB, ut utils.Utils) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return dbtx.WithTxContext(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx, `INSERT INTO roles(id, name) VALUES(?, ?)`)
		if err != nil {
			return fmt.Errorf("gagal prepare query insert roles: %w", err)
		}
		defer stmt.Close()

		roles := []string{"super admin", "admin", "editor", "penulis", "tamu"}
		for _, r := range roles {
			_, err := stmt.ExecContext(ctx, ut.GenerateULID(), r)
			if err != nil {
				return fmt.Errorf("gagal insert roles: %w", err)
			}
		}
		return nil
	})
}
