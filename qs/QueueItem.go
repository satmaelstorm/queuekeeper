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

func (qi QueueItem) String() string {
	return qi.message
}

func (qi QueueItem) Next() *QueueItem {
	return qi.next
}

func (qi QueueItem) LastAccess() int64 {
	return qi.lastAccess
}
