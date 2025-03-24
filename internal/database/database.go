package database

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// users table data
type User struct {
	ID        int
	Username  string
	Pass      string
	Email     string
	CreatedAt string
}

// uploads table data
type Upload struct {
	ID        int
	Filename  string
	Username  string
	CreatedAt string
	Enable    bool
}

// database
var DB *pgxpool.Pool

// initialize database connection
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

// close database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}

// get user data by username from users table
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

// update enable of a file
func UpdateUploadedFileInfoByID(id int, enable bool) error {
	query := `UPDATE uploads SET enable = $1 WHERE id=$2`

	if _, err := DB.Query(context.Background(), query, enable, id); err != nil {
		return err
	}
	return nil
}

// add a upload info
func AddUploadedFileInfo(filename, username string) error {
	query := `INSERT INTO uploads(filename, username) VALUES ($1, $2)`

	if _, err := DB.Query(context.Background(), query, filename, username); err != nil {
		return err
	}
	return nil
}

// TODO: Paginate data
// retrieve data on uploads table
func GetUploadsByUsername(username string) ([]Upload, error) {
	query := `SELECT id, filename, username, enable FROM uploads WHERE username = $1`

	rows, err := DB.Query(context.Background(), query, username)
	if err != nil {
		return nil, err
	}

	var uploads []Upload

	// TODO: better implementation with CollectRows
	for rows.Next() {
		var upload Upload
		if err := rows.Scan(&upload.ID, &upload.Filename, &upload.Username, &upload.Enable); err != nil {
			return nil, err
		}
		uploads = append(uploads, upload)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("error during row iteration")
	}

	return uploads, nil
}
