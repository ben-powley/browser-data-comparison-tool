package utils

import (
	"os"
	"path"
	"path/filepath"
)

func GetFilenamesFromFolder(folder string) ([]string, error) {
	workingDir, workingDirErr := os.Getwd()

	if workingDirErr != nil {
		return []string{}, workingDirErr
	}

	var filenames []string

	rootPath := path.Join(workingDir, folder)

	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil, err
	}

	pathWalkError := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".csv" {
			return nil
		}

		filenames = append(filenames, path)

		return nil
	})

	if pathWalkError != nil {
		return []string{}, pathWalkError
	}

	return filenames, nil
}
