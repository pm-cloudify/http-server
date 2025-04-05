package services

import (
	"errors"
	"log"
	"mime/multipart"

	"github.com/pm-cloudify/http-server/internal/api/v1/models"
	"github.com/pm-cloudify/http-server/internal/config"
	"github.com/pm-cloudify/shared-libs/acs3"
	database "github.com/pm-cloudify/shared-libs/psql"
)

// upload file to s3 and register in database
func UploadFile(file *multipart.FileHeader, username string) error {
	if file == nil || username == "" || file.Filename == "" {
		return errors.New("null args")
	}

	// save data to s3
	fileKey, err := acs3.UploadObject(config.AppConfigs.AC_S3Bucket, file)
	if err != nil {
		return errors.New("cannot save file")
	}

	// save record to db hold id of created record
	if err := database.AddUploadedFileInfo(file.Filename, username, fileKey); err != nil {
		err2 := acs3.DeleteObject(config.AppConfigs.AC_S3Bucket, fileKey)
		if err2 != nil {
			log.Printf("error: file with no record exists. key ->\"%s\"\n", fileKey)
		}
		log.Println(err.Error())
		return errors.New("failed to save record")
	}

	return nil
}

// returns list of saved files
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

	// TODO: find better implementation for parsing received data (or limit queries information / paginate)
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
