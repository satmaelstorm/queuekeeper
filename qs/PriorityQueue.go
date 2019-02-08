package qs

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



