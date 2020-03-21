package helpers

import (
	"io/ioutil"
	"os"
	"path"
)

func GetPath(filePath string) string {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return path.Join(workingDir, filePath)
}

func ReadFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return content
}
