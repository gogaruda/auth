package seeder

import (
	"database/sql"
	"github.com/gogaruda/auth/pkg/utils"
)

func SeedRun(db *sql.DB, ut utils.Utils) error {
	if err := Roles(db, ut); err != nil {
		return err
	}

	if err := Users(db, ut); err != nil {
		return err
	}

	return nil
}
