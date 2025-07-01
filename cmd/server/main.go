package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/auth/internal"
	"github.com/gogaruda/auth/internal/bootstrap"
	"github.com/gogaruda/auth/internal/config"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadConfig()
	db, err := bootstrap.SetupDatabase(cfg.DB)
	if err != nil {
		log.Fatal("koneksi database gagal:", err)
	}

	gin.SetMode(cfg.Mode.Debug)
	r := gin.Default()

	app := bootstrap.InitBootstrap(db, cfg)
	internal.RouteRegister(r, app)

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("server gagal dijalankan...")
	}
}
