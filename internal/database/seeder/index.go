package seeder

import "sql/pkg/system/config"

func SeedRun() error {
	db := config.DB
	if err := UserSeed(db); err != nil {
		return err
	}

	return nil
}
