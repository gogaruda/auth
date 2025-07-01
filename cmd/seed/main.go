package main

import (
	"fmt"
	"github.com/gogaruda/auth/database/seeders"
	"github.com/gogaruda/auth/internal/bootstrap"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/pkg/utils"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadConfig()
	db, err := bootstrap.SetupDatabase(cfg.DB)
	if err != nil {
		log.Fatal("koneksi ke database gagal:", err)
	}

	newID := utils.NewULIDCreate()
	hash := utils.NewBcryptHasher()

	if err := seeder.SeedRun(db, newID, hash); err != nil {
		log.Fatalf("seeding gagal: %w", err)
	}

	fmt.Print("seeder berhasil...")
}
