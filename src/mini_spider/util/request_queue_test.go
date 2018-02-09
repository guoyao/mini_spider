package util

import (
	"net/http"
	"strconv"
	"testing"
)

func TestPush(t *testing.T) {
	funcName := "Push"

	queue := NewRequestQueue()
	req, err := http.NewRequest("GET", "https://www.baidu.com", nil)
	if err != nil {
		t.Error(FormatTest(funcName, err.Error(), "nil"))
		return
	}

	request := NewRequest(req, 1)
	queue.Push(request)
	expected := 1
	result := queue.Len()

	if result != expected {
		t.Error(FormatTest(funcName, strconv.Itoa(result), strconv.Itoa(expected)))
	}
}

func TestPushAll(t *testing.T) {
	funcName := "PushAll"
	requests := make([]*Request, 0, 10)

	queue := NewRequestQueue()
	req, err := http.NewRequest("GET", "https://www.baidu.com", nil)
	if err != nil {
		t.Error(FormatTest(funcName, err.Error(), "nil"))
		return
	}
	request := NewRequest(req, 1)
	requests = append(requests, request)

	req, err = http.NewRequest("GET", "https://www.google.com", nil)
	if err != nil {
		t.Error(FormatTest(funcName, err.Error(), "nil"))
		return
	}
	request = NewRequest(req, 1)
	requests = append(requests, request)
	queue.PushAll(requests)

	expected := 2
	result := queue.Len()

	if result != expected {
		t.Error(FormatTest(funcName, strconv.Itoa(result), strconv.Itoa(expected)))
	}
}

func TestPop(t *testing.T) {
	funcName := "Pop"

	queue := NewRequestQueue()
	req, err := http.NewRequest("GET", "https://www.baidu.com", nil)
	if err != nil {
		t.Error(FormatTest(funcName, err.Error(), "nil"))
		return
	}

	request := NewRequest(req, 1)
	queue.Push(request)
	expected := 1
	result := queue.Len()

	if result != expected {
		t.Error(FormatTest(funcName, strconv.Itoa(result), strconv.Itoa(expected)))
	}

	queue.Pop()
	expected = 0
	result = queue.Len()

	if result != expected {
		t.Error(FormatTest(funcName, strconv.Itoa(result), strconv.Itoa(expected)))
	}
}
