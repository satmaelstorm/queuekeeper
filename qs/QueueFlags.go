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

type QueueFlagsSetter func(bool) *QueueFlags

/**
 * Getters
 */
func NewQueueFlags() QueueFlags {
	flags := QueueFlags{durable: false, deduplicated: false, delayedDelivery: 0, withPriority: false, stack: false}
	return flags
}

func (qf *QueueFlags) isDurable() bool {
	return qf.durable
}

func (qf *QueueFlags) isDeduplicated() bool {
	return qf.deduplicated
}

func (qf *QueueFlags) isDelayedDelivery() bool {
	return qf.delayedDelivery > 0
}

func (qf *QueueFlags) DelayedDelivery() int {
	return qf.delayedDelivery
}

func (qf *QueueFlags) isWithPriority() bool {
	return qf.withPriority
}

func (qf *QueueFlags) isStack() bool {
	return qf.stack
}

/**
 * Setters
 */
func (qf *QueueFlags) SetDurable(val bool) *QueueFlags {
	qf.durable = val
	return qf
}

func (qf *QueueFlags) SetDeduplicated(val bool) *QueueFlags {
	qf.deduplicated = val
	return qf
}

func (qf *QueueFlags) SetDelayedDelivery(val int) *QueueFlags {
	qf.delayedDelivery = val
	return qf
}

func (qf *QueueFlags) SetWithPriority(val bool) *QueueFlags {
	qf.withPriority = val
	return qf
}

func (qf *QueueFlags) SetStack(val bool) *QueueFlags {
	qf.stack = val
	return qf
}
