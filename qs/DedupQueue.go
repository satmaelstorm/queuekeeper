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

func (this *DedupQueue) Put(msg string) (*QueueItem, error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.put(msg)
}

func (this *DedupQueue) put(msg string) (*QueueItem, error) {
	qi, ok := this.items[msg]
	if ok {
		qi.lastAccess = getTimeForLastAccess()
		this.cnt += 1
		return qi, nil
	}

	qi, err := this.Queue.put(msg)

	if err != nil {
		return nil, err
	}

	this.items[msg] = qi
	this.cnt += 1

	return qi, nil
}
