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
	"mini_spider/util"
)

func TestFetch(t *testing.T) {
	funcName := "Fetch"
	fetcher := NewWebpageFetcher(5, ".")
	req, err := http.NewRequest("GET", "https://www.baidu.com", nil)
	if err != nil {
		t.Error(util.FormatTest(funcName, err.Error(), "nil"))
		return
	}

	media, err := fetcher.Fetch(req)
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
	fetcher := NewWebpageFetcher(5, outputDir)
	webpage := media.NewWebpage("http://test", bytes.NewReader(content), "", "")
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
