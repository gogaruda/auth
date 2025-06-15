package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadENV() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("File .env tidak ditemukan")
	}
}
