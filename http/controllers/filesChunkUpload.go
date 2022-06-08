package controllers

import (
	"cloudgobrrr/backend/pkg/helpers"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HttpFileUpload(c *gin.Context) {
	fileChunk, _, err := c.Request.FormFile("chunk")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "content-Type should be multipart/form-data"})
		return
	}

	fileName := c.PostForm("name")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	tempDir := filepath.Join(os.Getenv("TEMP_DIRECTORY"), c.MustGet("userName").(string))
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	rangeStart, rangeEnd, fileSize, err := helpers.GetContentRange(c.Request.Header.Get("Content-Range"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if err := helpers.ChunkUploadTmpMetaFile(rangeStart, rangeEnd, tempDir, fileName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		err := os.Mkdir(tempDir, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "error creating temporary directory"})
			return
		}
	}

	filePath := filepath.Join(tempDir, fileName)

	tempFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "error creating file"})
		return
	}

	if _, err := io.Copy(tempFile, fileChunk); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "error writing to a file"})
		return
	}
	tempFile.Close()

	if rangeEnd >= fileSize-1 {
		uploadingFile, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "failed to upload file"})
			return
		}
		uploadingFile.Close()

		c.JSON(http.StatusOK, gin.H{"status": "uploaded"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "uploading"})
}

type bindingFileUploadFinish struct {
	Path string `json:"path" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func HttpFileUploadFinish(c *gin.Context) {
	var json bindingFileUploadFinish
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	tempDir := os.Getenv("TEMP_DIRECTORY") + "/" + c.MustGet("userName").(string)

	tempFilePath := filepath.Join(tempDir, json.Name)

	tempFile, err := os.Open(tempFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "failed to open temp file"})
		return
	}
	defer tempFile.Close()

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), json.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	err = os.Rename(tempFilePath, filepath.Join(path, json.Name))
	if err != nil {
		fmt.Println("move error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "error moving file"})
		return
	}

	// delete meta file
	err = os.Remove(tempFilePath + ".meta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "error deleting meta file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "finished"})
}

// ToDo: reduce overhead
