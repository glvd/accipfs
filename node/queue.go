package node

import "time"

// Queue ...
type Queue struct {
	exchange *Exchange
	option   *QueueOption
	callback chan *Exchange
}

// QueueOption ...
type QueueOption struct {
	Callback bool
	Timeout  time.Duration
}

// CallbackWaiter ...
type CallbackWaiter interface {
	WaitCallback() *Exchange
}

// QueueOptions ...
type QueueOptions func(option *QueueOption)

// NewQueue ...
func NewQueue(exchange *Exchange, opts ...QueueOptions) *Queue {
	op := defaultQueueOption()
	for _, opt := range opts {
		opt(op)
	}

	q := &Queue{
		exchange: exchange,
		option:   op,
	}
	if q.option.Callback {
		q.callback = make(chan *Exchange)
	}
	return q
}

// HasCallback ...
func (q *Queue) HasCallback() bool {
	return q.option.Callback
}

// Exchange ...
func (q *Queue) Exchange() *Exchange {
	return q.exchange
}

// SetSession ...
func (q *Queue) SetSession(s uint32) {
	q.exchange.Session = s
}

// Callback ...
func (q *Queue) Callback(exchange *Exchange) {
	if q.callback != nil {
		t := time.NewTimer(q.option.Timeout)
		select {
		case <-t.C:
		case q.callback <- exchange:
			t.Reset(0)
		}
	}
}

// WaitCallback ...
func (q *Queue) WaitCallback() *Exchange {
	if q.callback != nil {
		t := time.NewTimer(q.option.Timeout)
		select {
		case <-t.C:
		case cb := <-q.callback:
			t.Reset(0)
			return cb
		}
	}
	return nil
}

// Send ...
func (q *Queue) Send(out chan<- *Queue) bool {
	if out == nil {
		return false
	}
	t := time.NewTimer(q.option.Timeout)
	select {
	case <-t.C:
		return false
	case out <- q:
		t.Reset(0)
		return true
	}
}
