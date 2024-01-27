package utils

import (
	"errors"
	"sync"
)

type QueueEasy struct {
	queue []byte
	mux   sync.Mutex
}

func InitQueueEasy() *QueueEasy {
	q := &QueueEasy{
		queue: make([]byte, 0),
	}
	return q
}

func (q *QueueEasy) PushArray(ba []byte) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.queue = append(q.queue, ba...)
}

func (q *QueueEasy) Push(b byte) {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.queue = append(q.queue, b)
}

func (q *QueueEasy) Pop() (byte, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	v, err := q.Pick()

	if err != nil {
		return 0, err
	}

	q.queue = q.queue[1:]
	return v, nil
}

func (q *QueueEasy) Pick() (byte, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	if len(q.queue) <= 0 {
		return 0, errors.New("Queue empty!")
	}

	return q.queue[0], nil
}

func (q *QueueEasy) PickArray() ([]byte, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	var qlen int = len(q.queue)
	if qlen <= 0 {
		return nil, errors.New("Queue empty!")
	}

	return q.queue[0:], nil
}

func (q *QueueEasy) PopArray(popqlen int) ([]byte, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	var qlen int = len(q.queue)

	if popqlen > qlen {
		return nil, errors.New("Queue too small")
	}

	if qlen <= 0 {
		return nil, errors.New("Queue empty!")
	}
	resultq := q.queue[0:qlen]
	q.queue = q.queue[popqlen:]

	return resultq, nil
}
