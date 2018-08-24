package fetcher

import (
	"bytes"
	"fmt"
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

func (h *HttpFetcher) Exist(media media.Media) bool {
	return h.storageDriver.Exist(media)
}

func (h *HttpFetcher) Fetch(req *http.Request, needBody func(media.Media) bool) (result media.Media, err error) {
	req.Header.Set("User-Agent", BROWSER_USER_AGENT)

	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	defer func() {
		if result != nil && needBody != nil && needBody(result) {
			buf, readBodyError := ioutil.ReadAll(resp.Body)
			if readBodyError != nil {
				result, err = nil, readBodyError
			} else {
				result.SetContent(bytes.NewReader(buf))
			}
		}
	}()

	if resp.StatusCode >= 400 {
		err = fmt.Errorf("status code is %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	contentLength, atoiError := strconv.Atoi(resp.Header.Get("Content-Length"))
	if atoiError != nil {
		contentLength = 0
	}

	charset := util.GetCharsetFromContentType(contentType)

	result = media.NewWebpage(req.URL.String(), nil, contentType, uint(contentLength), charset)
	return
}

func (h *HttpFetcher) Save(media media.Media) error {
	return h.storageDriver.Save(media)
}
