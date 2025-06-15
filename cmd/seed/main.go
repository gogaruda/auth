package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"sql/internal/database/seeder"
	"sql/pkg/system/config"
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
