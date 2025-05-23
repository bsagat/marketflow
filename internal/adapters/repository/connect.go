package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"log"
)

type PostgresDatabase struct {
	Db *sql.DB
}

func ConnectDB() *PostgresDatabase {
	log.Println("Starting database connection...")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect Database %s", err.Error())
	}

	// Sending Ping message
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to send ping message to the Database %s", err.Error())
	}
	return &PostgresDatabase{Db: db}
}
