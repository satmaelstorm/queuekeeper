package qs

import "sync"

type ICommonQueue interface {
	Put(qi *QueueItem) (*QueueItem, error)
	Get() (*QueueItem, error)
	Count() int
}

type Queue struct {
	head  *QueueItem
	tail  *QueueItem
	cnt   int
	flags QueueFlags
	mu    sync.Mutex
}

func NewQueue(fl QueueFlags) *Queue {
	return &Queue{head: nil, tail: nil, cnt: 0, flags: fl}
}

func (q *Queue) Count() int {
	return q.cnt
}

/**
 * Put message to queue
 */
func (q *Queue) Put(qi *QueueItem) (*QueueItem, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.put(qi)
}

/**
 * Get message from queue
 */
func (q *Queue) Get() (*QueueItem, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.get()
}

func (q *Queue) put(qi *QueueItem) (*QueueItem, error) {
	if nil == q.head {
		q.head = qi
	}

	if nil != q.tail {
		q.tail.next = qi
	}

	q.tail = qi

	q.cnt += 1

	return qi, nil
}

func (q *Queue) get() (*QueueItem, error) {
	if nil == q.head {
		return nil, NewError("head is nil", ErrQueueEmpty)
	}

	if q.flags.isDelayedDelivery() {
		return findMessageCanDelivery(q)
	}

	if 1 == q.cnt {
		q.tail = nil
	}
	qi := q.head
	q.head = qi.next
	q.cnt -= 1
	return qi, nil
}

func findMessageCanDelivery(q *Queue) (*QueueItem, error) {
	t := getTimeForLastAccess()
	curItem := q.head
	delay := curItem.delay
	if delay < 0 {
		delay = int64(q.flags.DelayedDelivery())
	}
	for (curItem.lastAccess + delay) >= t {
		if nil == curItem.next {
			return nil, NewError("no messages for delivery", ErrNothingToDelivery)
		}
		curItem = curItem.next
	}

	return curItem, nil
}
