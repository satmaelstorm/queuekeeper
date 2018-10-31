package qs

type QueueFlags struct {
	durable         bool
	deduplicated    bool
	delayedDelivery bool
	withPriority    bool
	stack           bool
}

func NewQueueFlags(durable bool, dedup bool, delayedDelivery bool, withPriority bool, stack bool) NewQueueFlags {
	flags = QueueFlags{durable: durable, deduplicated: dedup, delayedDelivery: delayedDelivery, withPriority: withPriority, stack: stack}
	return flags
}
