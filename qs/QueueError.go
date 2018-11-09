package qs

const (
	ErrGeneral           = 1
	ErrQueueEmpty        = 2
	ErrNothingToDelivery = 3
	ErrNoQueue           = 4
)

type QueueError struct {
	msg  string
	code int
}

func (this QueueError) Error() string {
	return this.msg
}

func (this QueueError) Code() int {
	return this.code
}

func NewError(m string, c int) QueueError {
	return QueueError{msg: m, code: c}
}
