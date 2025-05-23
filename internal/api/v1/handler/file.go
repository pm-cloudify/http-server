package handler

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pm-cloudify/http-server/internal/api/v1/models"
	"github.com/pm-cloudify/http-server/internal/api/v1/services"
)

const MaxFileSize = 1024 * 1024 // 1MB

// upload a file
func Upload(c *gin.Context) {
	data_type := strings.Split(c.Request.Header.Get("Content-Type"), ";")[0]
	if data_type != "multipart/form-data" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body data should be multipart/form-data"})
		return
	}

	// limiting file size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize)

	var data models.ReceivedFileRequest

	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong body format"})
		return
	}

	// checking if any file exists
	if data.File == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "empty file"})
		return
	}

	// validate files
	allowedTypes := map[string]bool{
		".py": true,
	}
	ext := filepath.Ext(data.File.Filename)
	if !allowedTypes[ext] {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid file type"})
		return
	}

	if c.GetString("username") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}

	if err := services.UploadFile(data.File, c.GetString("username")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file uploaded"})
}

// TODO: paginate
// get list of uploaded files
func GetListOfFiles(c *gin.Context) {
	// request for list of files
	if c.GetString("username") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized user"})
		return
	}

	data, err := services.GetListOfUploads(c.GetString("username"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot request for file info"})
		return
	}
	c.JSON(http.StatusOK, data)
}
