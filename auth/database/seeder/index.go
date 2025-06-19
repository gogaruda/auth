package seeder

import (
	"github.com/gogaruda/auth/auth/config"
)

func SeedRun() error {
	db := config.ConnectDB()

	if err := User_roles(db); err != nil {
		return err
	}

	return nil
}
