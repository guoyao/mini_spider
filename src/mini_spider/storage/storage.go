package storage

import (
	"mini_spider/log"
	"mini_spider/media"
	"mini_spider/util"
)

type StorageDriver interface {
	Exist(media media.Media) bool
	Save(media media.Media) error
}

func getFileName(media media.Media) string {
	name := media.Name()
	contentType := media.ContentType()

	if contentType == "application/pdf" {
		url, err := util.URLDecode(name)
		if err == nil {
			url, err = util.URLDecode(url)
		}

		if err != nil {
			log.Logger.Warn("URLDecode failed: " + err.Error())
			return name
		}

		fileName := util.FileNameFromUrl(url)
		if fileName == "" {
			return name
		}

		return fileName
	}

	return name
}
