package seeder

func SeedRun() error {
	if err := UserSeed(); err != nil {
		return err
	}

	return nil
}
