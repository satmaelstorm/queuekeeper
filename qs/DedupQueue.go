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

func (this *DedupQueue) Put(qi *QueueItem) (*QueueItem, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.put(qi)
}

func (this *DedupQueue) put(pqi *QueueItem) (*QueueItem, error) {
	qi, ok := this.items[pqi.message]
	if ok {
		qi.lastAccess = getTimeForLastAccess()
		qi.delay = pqi.delay
		this.cnt += 1
		return qi, nil
	}

	qi, err := this.Queue.put(pqi)

	if err != nil {
		return nil, err
	}

	this.items[qi.message] = qi
	this.cnt += 1

	return qi, nil
}
