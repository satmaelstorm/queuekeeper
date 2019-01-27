package qs

type DedupQueue struct {
	Queue
	items map[string]*QueueItem
}

func NewDedupQueue(fl QueueFlags) *DedupQueue {
	if !fl.isDeduplicated() {
		fl.SetDeduplicated(true)
	}
	return &DedupQueue{
		Queue: Queue{head: nil, tail: nil, cnt: 0, flags: fl},
		items: make(map[string]*QueueItem)}
}

func (dq *DedupQueue) Put(qi *QueueItem) (*QueueItem, error) {
	dq.mu.Lock()
	defer dq.mu.Unlock()
	return dq.put(qi)
}

func (dq *DedupQueue) put(pqi *QueueItem) (*QueueItem, error) {
	qi, ok := dq.items[pqi.message]
	if ok {
		qi.lastAccess = getTimeForLastAccess()
		qi.delay = pqi.delay
		dq.cnt += 1
		return qi, nil
	}

	qi, err := dq.Queue.put(pqi)

	if err != nil {
		return nil, err
	}

	dq.items[qi.message] = qi
	dq.cnt += 1

	return qi, nil
}
