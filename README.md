# README.md
## Install
```
go get github.com/gogaruda/auth@v1.2.0
```
## Penggunaan
```go
config.LoadENV()
if os.Getenv("GIN_MODE") == "release" {
	gin.SetMode(gin.ReleaseMode)
}

db := config.ConnectDB()
app := auth.InitAuthModule(db)

r := gin.Default()
api := r.Group("/api")
auth.RegisterAuthRoutes(api.Group("/auth"), app.AuthService, app.UserService)

port := os.Getenv("APP_PORT")
fmt.Println(port)
if port == "" {
	port = "8080"
}

_ = r.Run(":" + port)
```

## Migration and Seed
### Jalankan Otomatis
#### Migration
```go
// cmd/migrate/main.go
package main

import (
	"log"
	"github.com/joho/godotenv"

	authdb "github.com/gogaruda/auth/auth/database"
	blogdb "github.com/gogaruda/blog/blog/database"
)

func main() {
	_ = godotenv.Load()

	log.Println("🚀 Migrasi modul auth...")
	if err := authdb.RunMigration(); err != nil {
		log.Fatalf("❌ Gagal migrasi auth: %v", err)
	}

	log.Println("🚀 Migrasi modul blog...")
	if err := blogdb.RunMigration(); err != nil {
		log.Fatalf("❌ Gagal migrasi blog: %v", err)
	}

	log.Println("✅ Semua migrasi selesai")
}
```

#### Seeder
```go
// cmd/seed/main.go
package main

import (
	"log"
	"github.com/joho/godotenv"

	authseed "github.com/gogaruda/auth/auth/database/seeder"
	blogseed "github.com/gogaruda/blog/blog/database/seeder"
	"github.com/gogaruda/blog/blog/config"
)

func main() {
	_ = godotenv.Load()
	config.ConnectDB()

	log.Println("🚀 Seeder modul auth...")
	if err := authseed.SeedRun(); err != nil {
		log.Fatalf("❌ Gagal seeder auth: %v", err)
	}

	log.Println("🚀 Seeder modul blog...")
	if err := blogseed.SeedRun(); err != nil {
		log.Fatalf("❌ Gagal seeder blog: %v", err)
	}

	log.Println("✅ Semua seeder selesai")
}
```
### Buat File Migrasi
```
migrate create -ext sql -dir auth/database/migrations -seq create_users_table
```
#### Note:
 - ganti auth sesuaikan dengan directory project
### Jalankan Migrasi
```
go run ./cmd/migrate/main.go
```

### Buat File Seeder
```
go run cmd/seed/create.go product
```

### Jalankan Seeder
```
go run cmd/seed/main.go
```