package images

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveImage(alias string, meta multipart.FileHeader, file multipart.File) error {
	dir, err := os.Getwd()
	if err != nil {
		slog.Error("unexpected error when getting working directory", slog.String("error", err.Error()))
		return errors.New("unexpected error occured")
	}

	filename := meta.Filename
	if alias != "" {
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(meta.Filename))
	}

	fileLocation := filepath.Join(dir, "public/images", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		slog.Error("unexpected error occured when creating new file", slog.String("error", err.Error()))
		return errors.New("unexpected error occured")
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, file); err != nil {
		slog.Error("unexpected error occured when trying to copy image data", slog.String("error", err.Error()))
		return errors.New("unexpected error occured")
	}

	return nil
}
