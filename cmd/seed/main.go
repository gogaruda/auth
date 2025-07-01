package main

import (
	seeder "auth/database/seeders"
	"auth/internal/bootstrap"
	"auth/internal/config"
	"fmt"
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

	if err := seeder.SeedRun(db); err != nil {
		log.Fatalf("seeding gagal: %w", err)
	}

	fmt.Print("seeder berhasil...")
}
