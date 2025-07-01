package main

import (
	"auth/database"
	"auth/internal/bootstrap"
	"auth/internal/config"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadConfig()
	db, err := bootstrap.SetupDatabase(cfg.DB)
	if err != nil {
		log.Fatal("Koneksi ke database gagal:", err)
	}

	if err := database.RunMigration(db); err != nil {
		log.Fatalf("Migrasi gagal: %v", err)
	}
}
