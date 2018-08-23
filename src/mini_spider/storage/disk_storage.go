package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"mini_spider/media"
	"mini_spider/util"
)

type DiskStorage struct {
	outputDir string
}

func NewDiskStorage(outputDir string) *DiskStorage {
	return &DiskStorage{outputDir: outputDir}
}

func (d *DiskStorage) Exist(media media.Media) bool {
	if media.ContentLength() == 0 {
		if strings.HasPrefix(media.ContentType(), "text/") {
			return false
		}
	} else if media.ContentLength() < 50*1024 {
		return false
	}

	return util.CheckFileExists(filepath.Join(d.outputDir, getFileName(media)))
}

func (d *DiskStorage) Save(media media.Media) error {
	content := media.Content()
	if content == nil {
		return errors.New("content is nil")
	}

	fileName := getFileName(media)
	path := filepath.Join(d.outputDir, fileName)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	return file.Sync()
}
