package qs

import "time"

type QueueItem struct {
	next       *QueueItem
	message    string
	lastAccess int64
	delay      int64
}

func getTimeForLastAccess() int64 {
	return time.Now().Unix()
}

func NewQueueItem(msg string, d int64) *QueueItem {
	qi := &QueueItem{message: msg, next: nil, lastAccess: getTimeForLastAccess(), delay: d}
	return qi
}

func (this QueueItem) String() string {
	return this.message
}

func (this QueueItem) Next() *QueueItem {
	return this.next
}

func (this QueueItem) LastAccess() int64 {
	return this.lastAccess
}
