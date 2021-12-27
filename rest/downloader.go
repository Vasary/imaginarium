package rest

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"imaginarium/helper"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	cleanFilename := strings.Trim(strings.TrimSuffix(fileName, filepath.Ext(fileName)), " ")

	if len(cleanFilename) < 32 {
		return c.JSON(http.StatusBadRequest, Error{"Invalid file name"})
	}

	filePath := fmt.Sprintf("%s/%s", StoragePath, helper.SplitFileNameToPath(cleanFilename, ext))

	if !exists(filePath) {
		return c.JSON(http.StatusBadRequest, Error{"File is not exists"})
	}

	return c.File(filePath)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
