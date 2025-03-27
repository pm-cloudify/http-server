package services

import (
	"errors"

	"github.com/pm-cloudify/http-server/internal/database"
	"github.com/pm-cloudify/http-server/pkg/auth"
)

func Login(username, password string) (string, error) {
	var user *database.User
	user, err := database.GetUserByUsername(username)
	if err != nil {
		return "", errors.New("database failed")
	}

	if err == nil && user == nil {
		return "", errors.New("user not found")
	}

	if !auth.VerifyPassword(user.Pass, password) {
		return "", errors.New("wrong password")
	}

	token, err := auth.GenerateToken(user.Username)

	if err != nil {
		return "", err
	}

	return token, nil
}
