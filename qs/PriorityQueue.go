package qs

import (
	"math/rand"
	"time"
	)

type PriorityQueue struct {
	Queue
	sum int
	queueManager *QueueManager
}

func NewPriorityQueue(fl QueueFlags, qm *QueueManager) *PriorityQueue{
	if len(fl.withPriority) < 1 {
		panic("Priority Queue without priorities!")
	}

	sum := 0

	for _, w := range fl.withPriority {
		sum += w
	}

	return &PriorityQueue{
		Queue: Queue{head: nil, tail: nil, cnt: 0, flags: fl},
		sum: sum,
		queueManager:qm}
}



func (pq *PriorityQueue)Put(qi *QueueItem)(*QueueItem, error){
	return nil, NewError("You can't put to this queue.", ErrCantPutToPriorityQueue)
}

func (pq PriorityQueue)Get()(*QueueItem, error){
	queueName := pq.queueDecision()
	queue, err := pq.queueManager.GetQueue(queueName)
	if err == nil {
		return queue.Get()
	}
	return nil, err
}

func (pq PriorityQueue) queueDecision() string {
	curSum := 0
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := rnd.Intn(pq.sum) + 1
	for queueName, queueWeight := range pq.flags.withPriority {
		curSum += queueWeight
		if r <= curSum {
			return queueName
		}
	}
	return ""
}
