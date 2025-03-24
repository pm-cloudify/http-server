package services

import (
	"errors"
	"log"
	"mime/multipart"

	"github.com/pm-cloudify/http-server/internal/api/v1/models"
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

func GetListOfUploads(username string) (*models.FilesList, error) {
	if username == "" {
		return nil, errors.New("null args")
	}

	var gathered_data []database.Upload

	gathered_data, err := database.GetUploadsByUsername(username)
	if err != nil {
		log.Panicln(err.Error())
		return nil, err
	}

	// TODO: what is nilness
	// TODO: find better implementation for parsing received data (or limit queries information)
	var prepared_data []models.FileInfo
	for _, v := range gathered_data {
		var data models.FileInfo
		data.ID, data.Filename = v.ID, v.Filename
		prepared_data = append(prepared_data, data)
	}

	return &models.FilesList{
		Data: prepared_data,
	}, nil
}
