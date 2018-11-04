package qs

/**
 * Queue Settings Flags
 */
type QueueFlags struct {
	durable         bool
	deduplicated    bool
	delayedDelivery int
	withPriority    bool
	stack           bool
}

/**
 * Getters
 */
func NewQueueFlags() QueueFlags {
	flags := QueueFlags{durable: false, deduplicated: false, delayedDelivery: 0, withPriority: false, stack: false}
	return flags
}

func (this *QueueFlags) isDurable() bool {
	return this.durable
}

func (this *QueueFlags) isDeduplicated() bool {
	return this.deduplicated
}

func (this *QueueFlags) isDelayedDelivery() bool {
	return this.delayedDelivery > 0
}

func (this *QueueFlags) DelayedDelivery() int {
	return this.delayedDelivery
}

func (this *QueueFlags) isWithPriority() bool {
	return this.withPriority
}

func (this *QueueFlags) isStack() bool {
	return this.stack
}

/**
 * Setters
 */
func (this *QueueFlags) setDurable(val bool) *QueueFlags {
	this.durable = val
	return this
}

func (this *QueueFlags) setDeduplicated(val bool) *QueueFlags {
	this.deduplicated = val
	return this
}

func (this *QueueFlags) setDelayedDelivery(val int) *QueueFlags {
	this.delayedDelivery = val
	return this
}

func (this *QueueFlags) setWithPriority(val bool) *QueueFlags {
	this.withPriority = val
	return this
}

func (this *QueueFlags) setStack(val bool) *QueueFlags {
	this.stack = val
	return this
}
