package rest

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"imaginarium/helper"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
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

	err = c.Request().ParseMultipartForm(FileSize << 20)
	if err != nil {
		return c.String(http.StatusForbidden, "File size is too big")
	}

	for _, file := range files {
		result, err := uploadFile(file)

		if err != nil {
			fmt.Println(err)
		}

		images = append(images, result)
	}

	return c.JSON(http.StatusCreated, images)
}

func uploadFile(file *multipart.FileHeader) (map[string][]Image, error) {
	all := make(map[string][]Image)

	src, err := file.Open()

	if err != nil {
		return all, err
	}

	buff := make([]byte, 512)
	_, err = src.Read(buff)
	src.Seek(0, 0)
	fileType := http.DetectContentType(buff)

	if _, ok := SupportedTypes[fileType]; !ok {
		return all, errors.New("unsupported file type")
	}

	tmpFileName := random.String(32)
	ext := path.Ext(file.Filename)

	workDir := fmt.Sprintf("%s/%s", StoragePath, helper.GetDestinationDir(tmpFileName))
	helper.CreateDir(workDir)

	uploadedFilePath := fmt.Sprintf("%s/%s", StoragePath, helper.GetDestinationWithFile(tmpFileName, ext))

	dst, err := os.Create(uploadedFilePath)
	if err != nil {
		return all, err
	}

	if _, err = io.Copy(dst, src); err != nil {
		return all, err
	}

	src.Close()
	dst.Close()

	var results []Image
	results = append(results, Image{fmt.Sprintf("%s%s", helper.HashFileName(tmpFileName), ext), "original"})
	for key, context := range Contexts {
		img, _ := createResizedCopies(workDir, fmt.Sprintf("%s%s", helper.HashFileName(tmpFileName), ext), key, context)
		results = append(results, img)
	}

	all[file.Filename] = results
	return all, nil
}

func createResizedCopies(filePath string, fileName string, context string, size Size) (Image, error) {
	input, _ := os.Open(fmt.Sprintf("%s/%s", filePath, fileName))
	defer input.Close()

	ext := path.Ext(fileName)
	cleanedFileName := strings.TrimSuffix(fileName, ext)

	output, _ := os.Create(fmt.Sprintf("%s/%s_%s%s", filePath, cleanedFileName, context, path.Ext(fileName)))
	defer output.Close()

	input.Seek(0, 0)
	src, _, err := image.Decode(input)

	if err != nil {
		return Image{}, err
	}

	newImage := resize.Resize(size.Width, size.Height, src, resize.Lanczos3)
	jpeg.Encode(output, newImage, nil)

	return Image{fmt.Sprintf("%s_%s%s", cleanedFileName, context, ext), context}, nil
}
