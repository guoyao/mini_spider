package fetcher

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"mini_spider/media"
	"mini_spider/storage"
	"mini_spider/util"
)

func TestFetch(t *testing.T) {
	funcName := "Fetch"
	fetcher := NewHttpFetcher(5, storage.NewDiskStorage("."))
	req, err := http.NewRequest("GET", "https://www.baidu.com", nil)
	if err != nil {
		t.Error(util.FormatTest(funcName, err.Error(), "nil"))
		return
	}

	media, err := fetcher.Fetch(req, func(media media.Media) bool {
		return true
	})
	if err != nil {
		t.Error(util.FormatTest(funcName, err.Error(), "nil"))
		return
	}

	buf, err := ioutil.ReadAll(media.Content())
	if err != nil {
		t.Error(util.FormatTest(funcName, err.Error(), "nil"))
		return
	}

	if !strings.Contains(string(buf), "www.baidu.com") {
		t.Error(util.FormatTest(funcName, "wrong page content", "page content"))
	}
}

func TestSave(t *testing.T) {
	funcName, outputDir, content := "Save", ".", []byte("test content")
	fetcher := NewHttpFetcher(5, storage.NewDiskStorage(outputDir))
	webpage := media.NewWebpage("http://test", bytes.NewReader(content), "", 0, "")
	err := fetcher.Save(webpage)
	if err != nil {
		t.Error(util.FormatTest(funcName, err.Error(), "nil"))
		return
	}

	file, err := os.Open(filepath.Join(outputDir, webpage.Name()))
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()
	if err != nil {
		t.Error(util.FormatTest(funcName, err.Error(), "nil"))
		return
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(util.FormatTest(funcName, err.Error(), "nil"))
		return
	}

	expected := len(content)
	result := len(fileContent)
	if result != expected {
		t.Error(util.FormatTest(funcName, "file content length: "+strconv.Itoa(result),
			"response content lenght: "+strconv.Itoa(expected)))
	}
}
