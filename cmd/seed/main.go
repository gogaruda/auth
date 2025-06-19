package main

import (
	"fmt"
	"github.com/gogaruda/auth/internal/config"
	"github.com/gogaruda/auth/internal/database/seeder"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env tidak ditemukan")
	}

	config.ConnectDB()

	// Jalankan semua seeder
	if err := seeder.SeedRun(); err != nil {
		log.Fatalf("❌ Gagal seeding: %v", err)
	}

	fmt.Println("✅ Seeder selesai dijalankan")
}
