package storage

import (
	"mini_spider/media"

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

func (b *BosStorage) Save(media media.Media) error {
	objectKey := getFileName(media)
	_, err := b.client.PutObject(b.bucket, objectKey, media.Content(), nil, nil)
	return err
}
