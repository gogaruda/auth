package main

import (
	"fmt"
	"github.com/gogaruda/auth/auth/database"
	"github.com/joho/godotenv"
	"log"
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
