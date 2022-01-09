package rest

import (
	"encoding/base64"
	"errors"
	"github.com/labstack/echo/v4"
	"imaginarium/resizer"
	"imaginarium/uploader"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type Image struct {
	Name    string
	Context string
}

type Size struct {
	Width  uint
	Height uint
}

var Contexts map[string]Size

func Upload(c echo.Context) error {
	form, err := c.MultipartForm()

	if err != nil {
		return err
	}

	files := form.File["image"]
	var images []map[string][]Image

	for _, file := range files {
		result, err := uploadFile(file)

		if err != nil {
			return c.String(http.StatusBadRequest, "Unsupported file type")
		}

		images = append(images, result)
	}

	if len(images) == 0 {
		return c.String(http.StatusBadRequest, "There are no files were provided")
	}

	return c.JSON(http.StatusCreated, images)
}

func uploadFile(file *multipart.FileHeader) (map[string][]Image, error) {
	collection := make(map[string][]Image)

	src, err := file.Open()
	if err != nil {
		return collection, err
	}

	bytes := make([]byte, file.Size)
	_, err = src.Read(bytes)
	fileType := strings.Split(http.DetectContentType(bytes), ";")[0]

	if err := src.Close(); err != nil {
		return nil, err
	}

	if "text/plain" == fileType {
		decoded, err := base64.StdEncoding.DecodeString(string(bytes))
		if err != nil {
			panic(err)
		}

		fileType = http.DetectContentType(decoded)
		bytes = decoded
	}

	if _, ok := SupportedTypes[fileType]; !ok {
		return collection, errors.New("unsupported file type")
	}

	uploadedFile := uploader.Upload(bytes, file.Filename)

	var results []Image

	results = append(results, Image{filepath.Base(uploadedFile.Name()), "original"})
	for key, context := range Contexts {
		img := resizer.CreateResizedCopies(uploadedFile, key, context.Width, context.Height)
		results = append(results, Image{filepath.Base(img.Name()), key})
	}

	collection[file.Filename] = results
	return collection, nil
}
