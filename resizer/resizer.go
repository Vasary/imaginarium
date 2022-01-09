package resizer

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func CreateResizedCopies(file *os.File, context string, width uint, height uint) *os.File {
	var err error

	ext := path.Ext(file.Name())
	cleanedFileName := strings.TrimSuffix(filepath.Base(file.Name()), ext)

	input, _ := os.Open(file.Name())
	src, _, err := image.Decode(input)

	if err != nil {
		panic(err)
	}

	output, _ := os.Create(fmt.Sprintf("%s/%s_%s%s", filepath.Dir(file.Name()), cleanedFileName, context, path.Ext(file.Name())))

	newImage := resize.Resize(width, height, src, resize.Lanczos3)
	err = jpeg.Encode(output, newImage, nil)
	if err != nil {
		panic(err)
	}

	err = output.Close()
	if err != nil {
		panic(err)
	}

	return output
}
