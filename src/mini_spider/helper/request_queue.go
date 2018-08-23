package helper

import (
	"container/list"
	"errors"
	"mini_spider/util"
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

func (q *RequestQueue) Push(request *Request) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	key := util.URLEncode(request.URL.String())

	if _, ok := q.urlMap[key]; ok {
		return
	}

	q.list.PushBack(request)
	q.urlMap[key] = true
}

func (q *RequestQueue) PushAll(requests []*Request) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for _, request := range requests {
		key := util.URLEncode(request.URL.String())

		if _, ok := q.urlMap[key]; ok {
			continue
		}

		q.list.PushBack(request)
		q.urlMap[key] = true
	}
}

func (q *RequestQueue) Pop() (*Request, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.list.Len() <= 0 {
		return nil, nil
	}

	item := q.list.Front()
	q.list.Remove(item)

	if request, ok := item.Value.(*Request); ok {
		return request, nil
	}

	return nil, errors.New("queue item is not a valid request")
}

func (q *RequestQueue) Len() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return q.list.Len()
}
