package tools

import (
	"io"
	"os"
)

func MoveFile(src *os.File, dest string) (err error) {
	err = os.Rename(src.Name(), dest)
	return
}

func CopyFile(src *os.File, dest string) (file *os.File, err error) {
	destFile, err := os.Create(dest)
	if err != nil {
		return nil, err
	}
	defer func(destFile *os.File) {
		_ = destFile.Close()
	}(destFile)

	_, err = io.Copy(destFile, src)
	if err != nil {
		return nil, err
	}

	return destFile, nil
}
