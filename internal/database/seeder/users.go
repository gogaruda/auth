package seeder

import (
	"database/sql"
	"fmt"
)

func UserSeed(db *sql.DB) error {
	query := `
		INSERT INTO users (name, email) VALUES
		('Admin', 'admin@example.com'),
		('User', 'user@example.com')
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("users seeding failed: %w", err)
	}

	fmt.Println("✅ Users seeded")
	return nil
}
