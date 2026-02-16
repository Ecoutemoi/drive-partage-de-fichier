package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DB est la connexion globale à la base de données
var DB *sql.DB

// Init initialise la connexion MySQL
func Init() (*sql.DB, error) {
	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		user, pass, host, port, name,
	)

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetConnMaxLifetime(30 * time.Minute)

	// Test de connexion
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	DB = conn
	return DB, nil
}

// Close() ferme proprement la connexion DB
func Close() {
	if DB != nil {
		_ = DB.Close()
	}
}

// getEnv lit une variable d’environnement avec valeur par défaut
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
