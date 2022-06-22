package helpers

import (
	"cloudgobrrr/backend/pkg/structs"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func FilesList(path string) ([]structs.File, error) {
	fileReader, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer fileReader.Close()
	fileList, err := fileReader.Readdir(0)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	output := []structs.File{}

	for _, file := range fileList {
		Size := ""
		if !file.IsDir() {
			Size = ByteCountSI(file.Size())
		}
		Type := "file"
		if file.IsDir() {
			Type = "dir"
		}

		tmp := structs.File{Name: file.Name(), Type: Type, Size: Size, Modified: file.ModTime().Unix()}
		output = append(output, tmp)
	}
	// ToDo: add shared files (Type: "share")
	return output, nil
}

func GetAndCheckPath(username string, path string) (string, error) {
	dirPath := filepath.Join(os.Getenv("USER_DIRECTORY"), username, path)
	userPath := filepath.Join(os.Getenv("USER_DIRECTORY"), username)

	if !strings.HasPrefix(dirPath, userPath) {
		return "", fmt.Errorf("invalid path")
	}

	return dirPath, nil
}

func ChunkedUploadMetaFile(rangeStart int, rangeEnd int, tempDir string, fileName string) error {
	filePath := filepath.Join(tempDir, fileName)
	if rangeStart == 0 {
		os.Remove(filePath)
		os.Remove(filePath + ".meta")

		metaFile, err := os.Create(filePath + ".meta")
		if err != nil {
			return err
		}
		_, err = metaFile.WriteString(fmt.Sprintf("%d", rangeEnd))
		if err != nil {
			return err
		}
		metaFile.Close()
	} else {
		content, err := ioutil.ReadFile(filePath + ".meta")
		if err != nil {
			return err
		}
		fileSize, err := strconv.Atoi(string(content))
		if err != nil {
			return err
		}
		if rangeStart != fileSize+1 {
			return fmt.Errorf("invalid range")
		}

		// set new content range in meta file
		metaFile, err := os.OpenFile(tempDir+"/"+fileName+".meta", os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		_, err = metaFile.WriteString(fmt.Sprintf("%d", rangeEnd))
		if err != nil {
			return err
		}
		metaFile.Close()
	}
	return nil
}
