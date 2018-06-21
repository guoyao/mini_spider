package fetcher

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"mini_spider/log"
	"mini_spider/media"
	"mini_spider/util"
)

const (
	BROWSER_USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 " +
		"(KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36"
)

var httpClient *http.Client

type WebpageFetcher struct {
	timeout   uint
	outputDir string
}

func NewWebpageFetcher(timeout uint, outputDir string) *WebpageFetcher {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: time.Duration(timeout) * time.Second}
	}

	return &WebpageFetcher{timeout: timeout, outputDir: outputDir}
}

func (w *WebpageFetcher) Fetch(req *http.Request) (media.Media, error) {
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
	charset := util.GetCharsetFromContentType(contentType)

	return media.NewWebpage(req.URL.String(), bytes.NewReader(buf), contentType, charset), nil
}

func (w *WebpageFetcher) Save(media media.Media) error {
	content := media.Content()
	if content == nil {
		return errors.New("content is nil")
	}

	fileName := getFileName(media)
	path := filepath.Join(w.outputDir, fileName)
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
