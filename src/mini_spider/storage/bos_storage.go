package storage

import (
	"mini_spider/media"
	"strings"

	"github.com/guoyao/baidubce-sdk-go/bce"
	"github.com/guoyao/baidubce-sdk-go/bos"
)

type BosStorage struct {
	bucket string
	client *bos.Client
}

func NewBosStorage(ak, sk, bucket string) *BosStorage {
	credentials := bce.NewCredentials(ak, sk)
	config := bos.NewConfig(bce.NewConfig(credentials))
	return &BosStorage{bucket: bucket, client: bos.NewClient(config)}
}

func (b *BosStorage) Exist(media media.Media) bool {
	if media.ContentLength() == 0 {
		if strings.HasPrefix(media.ContentType(), "text/") {
			return false
		}
	} else if media.ContentLength() < 50*1024 {
		return false
	}

	objectKey := getFileName(media)
	_, err := b.client.GetObjectMetadata(b.bucket, objectKey, nil)
	return err == nil
}

func (b *BosStorage) Save(media media.Media) error {
	objectKey := getFileName(media)
	objectMetadata := &bos.ObjectMetadata{StorageClass: bos.STORAGE_CLASS_COLD}
	_, err := b.client.PutObject(b.bucket, objectKey, media.Content(), objectMetadata, nil)
	return err
}
