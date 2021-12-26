package rest

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"imaginarium/helper"
	"net/http"
	"os"
	"path"
)

type Error struct {
	Message string `json:"message"`
}

func Download(c echo.Context) error {
	fileName := c.Param("name")
	ext := path.Ext(c.Param("name"))

	if ext == "" {
		return c.JSON(http.StatusBadRequest, Error{"File extension is not provided"})
	}

	filePath := fmt.Sprintf("%s/%s", StoragePath, helper.SplitFileNameToPath(fileName[0:32], ext))

	if !exists(filePath) {
		return c.JSON(http.StatusBadRequest, Error{"File is not exists"})
	}

	return c.File(filePath)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
