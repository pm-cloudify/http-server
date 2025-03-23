package services

import (
	"errors"

	"github.com/pm-cloudify/http-server/internal/database"
)

func Login(username, password string) (string, error) {
	// TODO: add a better authentication, also pass should not be saved as plain text!
	var user *database.User
	user, err := database.GetUserByUsername(username)

	if err != nil {
		return "", errors.New("user not fount")
	}

	if user.Pass != password {
		return "", errors.New("wrong password")
	}

	return "test", nil
}
