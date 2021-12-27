package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
)

func GetDestinationDir(fileName string) string {
	encoder := md5.New()
	encoder.Write([]byte(fileName))
	tmp := hex.EncodeToString(encoder.Sum(nil))

	return fmt.Sprintf("%s/%s/%s", tmp[0:1], tmp[1:4], tmp[4:10])
}

func GetDestinationWithFile(fileName string, ext string) string {
	return SplitFileNameToPath(HashFileName(fileName), ext)
}

func HashFileName(fileName string) string {
	encoder := md5.New()
	encoder.Write([]byte(fileName))
	return hex.EncodeToString(encoder.Sum(nil))
}

func SplitFileNameToPath(fileName string, ext string) string {
	return fmt.Sprintf("%s/%s/%s/%s%s", fileName[0:1], fileName[1:4], fileName[4:10], fileName, ext)
}

func CreateDir(dirName string) error {
	err := os.MkdirAll(dirName, 0755)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
