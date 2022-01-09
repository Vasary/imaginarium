package uploader

import (
	"fmt"
	"github.com/labstack/gommon/random"
	"imaginarium/helper"
	"os"
	"path"
)

var StoragePath string

func Upload(buffer []byte, fileName string) *os.File {
	var err error

	name := getTempFileName()
	workDir := fmt.Sprintf("%s/%s", StoragePath, helper.GetDestinationDir(name))

	err = helper.CreateDir(workDir)
	if err != nil {
		panic(err)
	}

	ext := path.Ext(fileName)
	uploadedFilePath := fmt.Sprintf("%s/%s", StoragePath, helper.GetDestinationWithFile(name, ext))
	newFile, _ := os.Create(uploadedFilePath)

	if _, err := newFile.Write(buffer); err != nil {
		panic(err)
	}

	newFile.Close()

	return newFile
}

func getTempFileName() string {
	return random.String(32)
}
