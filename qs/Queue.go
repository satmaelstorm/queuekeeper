package qs

type ICommonQueue interface {
	Put(msg string) (*IQueueItem, error)
	Get() (*IQueueItem, error)
}

type Queue struct {
	head *QueueItem
	tail *QueueItem
	cnt  int
}

func NewQueue() *Queue {
	return &Queue{head: nil, tail: nil, cnt: 0}
}

func (this *Queue) Put(msg string) (*IQueueItem, error) {
	qi := &QueueItem{message: msg, next: nil}
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

func (this *Queue) Get() (*IQueueItem, error) {
	if nil == this.head {
		return nil, NewError("head is nil", ErrQueueEmpty)
	}
	if 1 == this.cnt {
		this.tail = nil
	}
	qi := this.head
	this.head = qi.next
	this.cnt -= 1
	return qi, nil
}
