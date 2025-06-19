# README.md
## Install
```
go get github.com/gogaruda/auth@v1.1.1
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

### Buat File Migrasi
```
migrate create -ext sql -dir internal/database/migrations -seq create_users_table
```

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