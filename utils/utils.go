package utils

import (
	"io/ioutil"
	"os"
)

func IsGitFolder(folder string) bool {
	if _, err := os.Stat(folder + "/.git"); err == nil {
		return true
	} else {
		return false
	}
}

func GetDirs(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			files = append(files, file.Name())
		}
	}
	return files, nil
}
