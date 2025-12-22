package main

import (
	"log"
	"os"
	"time"

	"handsoft/internal/db"
	"handsoft/internal/http/routes"
	"handsoft/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carga variables desde .env (si existe)
	if err := godotenv.Load(); err != nil {
		log.Println("warning: no se encontró .env (se usarán variables del sistema)")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL no está definido. Revisa backend/.env o tus variables de entorno")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-me"
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	gormDB, err := db.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Migraciones (usa la lista centralizada en internal/models/migrate.go)
	if err := gormDB.AutoMigrate(models.Models()...); err != nil {
		log.Fatal(err)
	}

	// Seed base (lo haremos ahora/enseguida si quieres)
	// if err := db.SeedBaseData(gormDB); err != nil {
	// 	log.Fatal(err)
	// }

	r := gin.Default()

	routes.Register(r, routes.Deps{
		DB:        gormDB,
		JWTSecret: jwtSecret,
		Issuer:    "handsoft-api",
		AccessTTL: 24 * time.Hour,
	})

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
