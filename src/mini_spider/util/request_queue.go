package util

import (
	"container/list"
	"errors"
	"sync"
)

type RequestQueue struct {
	mutex  *sync.Mutex
	list   *list.List
	urlMap map[string]bool
}

func NewRequestQueue() *RequestQueue {
	return &RequestQueue{mutex: &sync.Mutex{}, list: list.New(), urlMap: make(map[string]bool)}
}

func (queue *RequestQueue) Push(request *Request) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	key := URLEncode(request.URL.String())

	if _, ok := queue.urlMap[key]; ok {
		return
	}

	queue.list.PushBack(request)
	queue.urlMap[key] = true
}

func (queue *RequestQueue) PushAll(requests []*Request) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	for _, request := range requests {
		key := URLEncode(request.URL.String())

		if _, ok := queue.urlMap[key]; ok {
			continue
		}

		queue.list.PushBack(request)
		queue.urlMap[key] = true
	}
}

func (queue *RequestQueue) Pop() (*Request, error) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if queue.list.Len() <= 0 {
		return nil, nil
	}

	item := queue.list.Front()
	queue.list.Remove(item)

	if request, ok := item.Value.(*Request); ok {
		return request, nil
	}

	return nil, errors.New("queue item is not a valid request")
}

func (queue *RequestQueue) Len() int {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	return queue.list.Len()
}
