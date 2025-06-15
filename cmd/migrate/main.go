package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"sql/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env tidak ditemukan")
	}

	if err := database.RunMigration(); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	fmt.Println("✅ Migration applied successfully")
}
