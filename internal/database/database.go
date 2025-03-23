package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID       int
	Username string
	Pass     string
	Email    string
}

var DB *pgxpool.Pool

func InitDB(connectionStr string) error {
	config, err := pgxpool.ParseConfig(connectionStr)
	if err != nil {
		return err
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	log.Println("Successfully connected to the database")

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}

func GetUserByUsername(username string) (*User, error) {
	// TODO: check for query injection ?
	query := `SELECT id, username, pass, email FROM users WHERE username = $1`

	row := DB.QueryRow(context.Background(), query, username)

	var user User

	err := row.Scan(&user.ID, &user.Username, &user.Pass, &user.Email)

	if err != nil {
		log.Printf("Failed to fetch user: %s\n", err.Error())
		return nil, err
	}

	return &user, nil
}
