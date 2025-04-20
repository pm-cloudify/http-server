package models

import (
	"mime/multipart"
)

// login request body
// can be from a form or a json body
type LoginRequest struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

// upload file request body
// only accepting multipart/form-data
type ReceivedFileRequest struct {
	Inputs string                `form:"inputs"`
	File   *multipart.FileHeader `form:"file"`
}

// file info data
type FileInfo struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
}

// list of files information
type FilesList struct {
	Data []FileInfo `json:"data"`
}

// job request body
type JobRequest struct {
	FileID uint `form:"file_id" json:"file_id"`
}
