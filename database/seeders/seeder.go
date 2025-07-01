package seeder

import (
	"database/sql"
	"github.com/gogaruda/auth/pkg/utils"
)

func SeedRun(db *sql.DB, newID *utils.ULIDCreate, hash *utils.BcryptHasher) error {
	if err := Roles(db, newID); err != nil {
		return err
	}

	if err := Users(db, newID, hash); err != nil {
		return err
	}

	return nil
}
