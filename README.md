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