package fetcher

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"mini_spider/media"
)

const (
	BROWSER_USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36"
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

	return media.NewWebpage(req, resp), nil
}

func (w *WebpageFetcher) Save(media media.Media) error {
	content, err := media.GetContent()
	if err != nil {
		return err
	}

	path := filepath.Join(w.outputDir, media.GetName())
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)

	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return file.Sync()
}
