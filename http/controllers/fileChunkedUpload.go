package controllers

import (
	"cloudgobrrr/backend/http/binding"
	"cloudgobrrr/backend/pkg/helpers"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HttpFileChunkedUpload(c *gin.Context) {
	fileChunk, _, err := c.Request.FormFile("chunk")
	if err != nil {
		c.JSON(http.StatusBadRequest, binding.ResError{Error: "content-Type should be multipart/form-data"})
		return
	}

	fileName := c.PostForm("name")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	tempDir := filepath.Join(os.Getenv("TEMP_DIRECTORY"), c.MustGet("userName").(string))
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	rangeStart, rangeEnd, fileSize, err := helpers.HttpGetContentRange(c.Request.Header.Get("Content-Range"))
	if err != nil {
		c.JSON(http.StatusBadRequest, binding.ResError{Error: err.Error()})
		return
	}

	if err := helpers.ChunkedUploadMetaFile(rangeStart, rangeEnd, tempDir, fileName); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResError{Error: err.Error()})
		return
	}

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		err := os.Mkdir(tempDir, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
			return
		}
	}

	filePath := filepath.Join(tempDir, fileName)

	tempFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	if _, err := io.Copy(tempFile, fileChunk); err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}
	tempFile.Close()

	if rangeEnd >= fileSize-1 {
		uploadingFile, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
			return
		}
		uploadingFile.Close()

		c.JSON(http.StatusOK, binding.ResEmpty)
		return
	}
	c.JSON(http.StatusAccepted, binding.ResEmpty)
}

func HttpFileChunkedUploadFinish(c *gin.Context) {
	var json binding.ReqFileChunkedUploadFinish
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	tempDir := os.Getenv("TEMP_DIRECTORY") + "/" + c.MustGet("userName").(string)

	tempFilePath := filepath.Join(tempDir, json.Name)

	tempFile, err := os.Open(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}
	defer tempFile.Close()

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), json.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidPath)
		return
	}

	err = os.Rename(tempFilePath, filepath.Join(path, json.Name))
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	// delete meta file
	err = os.Remove(tempFilePath + ".meta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, binding.ResEmpty)
}

// ToDo: reduce overhead
