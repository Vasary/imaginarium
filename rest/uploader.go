package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"imaginarium/helper"
	"io"
	"net/http"
	"os"
	"path"
)

type Image struct {
	Original string
	Path     string
}

func Upload(c echo.Context) error {
	form, err := c.MultipartForm()

	if err != nil {
		return err
	}

	files := form.File["image"]
	var images []Image

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}

		err = c.Request().ParseMultipartForm(FileSize << 20)
		if err != nil {
			return c.String(http.StatusForbidden, "File size is too big")
		}

		buff := make([]byte, 512)
		_, err = src.Read(buff)
		src.Seek(0, 0)
		fileType := http.DetectContentType(buff)

		if _, ok := SupportedTypes[fileType]; !ok {
			return c.String(http.StatusForbidden, "Unsupported file type")
		}

		tmpFileName := random.String(32)
		ext := path.Ext(file.Filename)
		helper.CreateDir(fmt.Sprintf("%s/%s", StoragePath, helper.GetDestinationDir(tmpFileName)))

		dstFilePath := fmt.Sprintf("%s/%s", StoragePath, helper.GetDestinationWithFile(tmpFileName, ext))

		dst, err := os.Create(dstFilePath)
		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}

		if _, err = io.Copy(dst, src); err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}

		images = append(images, Image{file.Filename, fmt.Sprintf("%s%s", tmpFileName, ext)})

		src.Close()
		dst.Close()
	}

	return c.JSON(http.StatusCreated, images)
}
