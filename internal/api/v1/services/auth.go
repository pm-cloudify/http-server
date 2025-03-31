package services

import (
	"errors"
	"log"
	"regexp"

	"github.com/pm-cloudify/http-server/internal/database"
	"github.com/pm-cloudify/http-server/pkg/auth"
)

// var PasswordPattern = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]+$`)
var UsernamePattern = regexp.MustCompile(`^[a-zA-Z]{8,}$`)

func isValidUsername(username string) error {
	if len(username) < 4 {
		return errors.New("short username")
	}
	if len(username) > 64 {
		return errors.New("long username")
	}
	if !UsernamePattern.MatchString(username) {
		return errors.New("invalid username")
	}
	return nil
}

func checkPasswordPattern(password string) bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		validChars = true
	)

	for _, c := range password {
		switch {
		case 'a' <= c && c <= 'z':
			hasLower = true
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case '0' <= c && c <= '9':
			hasNumber = true
		default:
			validChars = false
		}
	}

	return hasLower && hasUpper && hasNumber && validChars
}

// TODO: use libs like validator
func isPasswordValid(password string) error {
	if len(password) < 8 {
		return errors.New("short password")
	}
	if len(password) > 256 {
		return errors.New("long password")
	}
	// if !PasswordPattern.MatchString(password) {
	if !checkPasswordPattern(password) {
		return errors.New("invalid password")
	}
	return nil
}

// login user. returns a token for user
func Login(username, password string) (string, error) {
	var user *database.User
	user, err := database.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("database failed")
	}

	if user == nil {
		return "", errors.New("user not found")
	}
	isPassValid, err := auth.VerifyPassword(user.Pass, password)
	if err != nil {
		return "", errors.New("verification failed")
	}
	if isPassValid {
		return "", errors.New("wrong password")
	}

	token, err := auth.GenerateToken(user.Username)

	if err != nil {
		return "", err
	}

	return token, nil
}

// TODO: add email for users
// creates a new user.
func SingIn(username, password string) error {

	// validate user and pass
	if err := isPasswordValid(password); err != nil {
		return errors.New("invalid password")
	}
	if err := isValidUsername(username); err != nil {
		return errors.New("invalid username")
	}

	// generate hashed password
	hashed_pass, err := auth.GenerateHash(password, auth.DefaultArgon2Params)
	if err != nil {
		log.Printf("hash-error: %s", err.Error())
		return errors.New("failed to create account")
	}

	// insert user to db
	err = database.AddUser(username, hashed_pass)
	if err != nil {
		log.Printf("db-error: %s", err.Error())
		return errors.New("failed to create account")
	}

	return nil
}
