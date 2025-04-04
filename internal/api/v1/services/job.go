package services

import (
	"errors"
	"fmt"

	"github.com/pm-cloudify/http-server/internal/config"
	"github.com/pm-cloudify/shared-libs/mb"
	database "github.com/pm-cloudify/shared-libs/psql"
)

func SendRunRequest(username string, file_id uint) error {

	// 1 - get file by file_id
	file_data, err := database.GetUploadByFileId(file_id)
	if file_data == nil && err == nil {
		return errors.New("file not found")
	}
	if err != nil {
		return errors.New("failed to get file data")
	}

	// 2 - check if username is the owner of the file
	if file_data.Username != username {
		return errors.New("file not owned by user")
	}

	// 3 - check if file is never sent
	if file_data.Enable {
		return errors.New("file already sent")
	}

	// 4 - create a message and send it to RMQ
	err = mb.ProduceTextMsg(config.App_MB, fmt.Sprintf("file_id=%d", file_id))
	if err != nil {
		return errors.New("failed to send request")
	}

	// 5 - update file data to enable = true
	err = database.UpdateUploadEnableByFileId(file_id, true)
	if err != nil {
		// TODO: hold incomplete updates in a temp cache to update db later
		// example: use redis or temp mem
		return errors.New("failed to update file data")
	}

	return nil
}
