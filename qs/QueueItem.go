package qs

type IQueueItem interface {
	Next() IQueueItem
	String() string
}

type QueueItem struct {
	next    *QueueItem
	message string
}

type PriorityQueueItem struct {
	QueueItem
	weight int
}

func (this QueueItem) String() string {
	return this.message
}

func (this QueueItem) Next() IQueueItem {
	return this.next
}
