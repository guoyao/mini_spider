package fetcher

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"mini_spider/media"
	"mini_spider/storage"
	"mini_spider/util"
)

const (
	BROWSER_USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36"
)

var httpClient *http.Client

type HttpFetcher struct {
	timeout       uint
	storageDriver storage.StorageDriver
}

func NewHttpFetcher(timeout uint, storageDriver storage.StorageDriver) *HttpFetcher {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	}

	return &HttpFetcher{timeout: timeout, storageDriver: storageDriver}
}

func (h *HttpFetcher) GetMetadata(r *http.Request) (media.Media, error) {
	url := r.URL.String()
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", BROWSER_USER_AGENT)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		contentLength = 0
	}

	charset := util.GetCharsetFromContentType(contentType)

	return media.NewWebpage(url, nil, contentType, uint(contentLength), charset), nil
}

func (h *HttpFetcher) Exist(metadata media.Media) bool {
	return h.storageDriver.Exist(metadata)
}

func (h *HttpFetcher) Fetch(req *http.Request) (media.Media, error) {
	req.Header.Set("User-Agent", BROWSER_USER_AGENT)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("status code is " + strconv.Itoa(resp.StatusCode))
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	contentType := resp.Header.Get("Content-Type")
	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		contentLength = 0
	}

	charset := util.GetCharsetFromContentType(contentType)

	return media.NewWebpage(req.URL.String(), bytes.NewReader(buf), contentType, uint(contentLength), charset), nil
}

func (h *HttpFetcher) Save(media media.Media) error {
	return h.storageDriver.Save(media)
}
