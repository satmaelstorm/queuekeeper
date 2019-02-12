package qs

const (
	ErrGeneral                = 1
	ErrQueueEmpty             = 2
	ErrNothingToDelivery      = 3
	ErrNoQueue                = 4
	ErrCantPutToPriorityQueue = 5
)

type QueueError struct {
	msg  string
	code int
}

func (qe QueueError) Error() string {
	return qe.msg
}

func (qe QueueError) Code() int {
	return qe.code
}

func NewError(m string, c int) QueueError {
	return QueueError{msg: m, code: c}
}
