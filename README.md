# README.md
## Install
```go
go get github.com/gogaruda/auth@v1.0.0
```
## Penggunaan
```go
config.LoadENV()
if os.Getenv("GIN_MODE") == "release" {
	gin.SetMode(gin.ReleaseMode)
}

r := gin.Default()

app := container.InitApp()
routes.InitRouter(r, app)

port := os.Getenv("APP_PORT")
fmt.Println(port)
if port == "" {
	port = "8080"
}

r.Run(":" + port)
```

### Buat File Migrasi
```go
migrate create -ext sql -dir internal/database/migrations -seq create_users_table
```

### Jalankan Migrasi
```cmd
go run ./cmd/migrate/main.go
```

### Buat File Seeder
```go
go run cmd/seed/create.go product
```

### Jalankan Seeder
```cmd
go run cmd/seed/main.go
```