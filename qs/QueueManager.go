package qs

import (
	"fmt"
	"strings"
)

type QueueManager struct {
	queues map[string]ICommonQueue
}

func NewQueueManager() *QueueManager {
	return &QueueManager{queues: make(map[string]ICommonQueue)}
}

func (this *QueueManager) CreateQueue(name string, flags QueueFlags) ICommonQueue {
	if flags.isDeduplicated() {
		this.queues[name] = NewDedupQueue(flags)
	} else {
		this.queues[name] = NewQueue(flags)
	}
	return this.queues[name]
}

func (this *QueueManager) GetQueue(name string) (ICommonQueue, error) {
	q, ok := this.queues[name]
	if ok {
		return q, nil
	}
	return nil, NewError("No Queue "+name, ErrNoQueue)
}

func (this *QueueManager) String() string {
	qStr := "Queues: "
	var keys []string
	for k := range this.queues {
		q := this.queues[k]
		item := k + fmt.Sprintf("(items: %d)", q.Count())
		keys = append(keys, item)
	}
	if len(keys) < 1 {
		return qStr + "NO QUEUES CONFIGURED"
	}
	return qStr + strings.Join(keys[:], "; ")

}
