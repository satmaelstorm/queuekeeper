package qs

import "sync"

type ICommonQueue interface {
	Put(qi *QueueItem) (*QueueItem, error)
	Get() (*QueueItem, error)
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

/**
 * Put message to queue
 */
func (this *Queue) Put(qi *QueueItem) (*QueueItem, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.put(qi)
}

/**
 * Get message from queue
 */
func (this *Queue) Get() (*QueueItem, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.get()
}

func (this *Queue) put(qi *QueueItem) (*QueueItem, error) {
	if nil == this.head {
		this.head = qi
	}

	if nil != this.tail {
		this.tail.next = qi
	}

	this.tail = qi

	this.cnt += 1

	return qi, nil
}

func (this *Queue) get() (*QueueItem, error) {
	if nil == this.head {
		return nil, NewError("head is nil", ErrQueueEmpty)
	}

	if this.flags.isDelayedDelivery() {
		return findMessageCanDelivery(this)
	}

	if 1 == this.cnt {
		this.tail = nil
	}
	qi := this.head
	this.head = qi.next
	this.cnt -= 1
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
