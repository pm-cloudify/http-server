package services

import (
	"errors"
	"mime/multipart"

	"github.com/pm-cloudify/http-server/internal/database"
)

// upload file to s3 and register in database
func UploadFile(file *multipart.FileHeader, username string) error {
	if file == nil || username == "" || file.Filename == "" {
		return errors.New("null args")
	}

	// save record to db. hold id of created record
	filename := file.Filename
	if err := database.AddUploadedFileInfo(filename, username); err != nil {
		return err
	}

	// TODO: 2 - upload file to s3 storage. user id of created record : save (id-filename.ext)
	return nil
}
