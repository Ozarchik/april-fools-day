package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func NewPostgres() (*sql.DB, error) {
	connectParams := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connectParams)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err == nil {
		fmt.Println("Connected to DB")
		return db, nil
	}

	return nil, fmt.Errorf("could not connect to DB: %w", err)
}
