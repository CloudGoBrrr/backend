package migrator

import (
	"cloudgobrrr/backend/pkg/env"
	"os"
	"path/filepath"
)

func RunMigrations() {
	versionFilePath := filepath.Join(os.Getenv("DATA_DIRECTORY"), "version.txt")

	// For future use maybe
	file, err := os.Create(versionFilePath)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(env.GetVersion())
	if err != nil {
		panic(err)
	}
	file.Close()

	// ToDo: write migration system
}
