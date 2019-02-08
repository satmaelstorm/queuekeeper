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

func (qm *QueueManager) CreateQueue(name string, flags QueueFlags) ICommonQueue {
	if len(flags.getWithPriority()) > 0 {
		qm.queues[name] = NewPriorityQueue(flags, qm)
	}else if flags.isDeduplicated() {
		qm.queues[name] = NewDedupQueue(flags)
	} else {
		qm.queues[name] = NewQueue(flags)
	}
	return qm.queues[name]
}

func (qm *QueueManager) GetQueue(name string) (ICommonQueue, error) {
	q, ok := qm.queues[name]
	if ok {
		return q, nil
	}
	return nil, NewError("No Queue "+name, ErrNoQueue)
}

func (qm *QueueManager) String() string {
	qStr := "Queues: "
	var keys []string
	for k := range qm.queues {
		q := qm.queues[k]
		item := k + fmt.Sprintf("(items: %d)", q.Count())
		keys = append(keys, item)
	}
	if len(keys) < 1 {
		return qStr + "NO QUEUES CONFIGURED"
	}
	return qStr + strings.Join(keys[:], "; ")

}
