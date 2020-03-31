package task

import (
	"context"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type task struct {
	MaxLimit int64
	stop     *atomic.Bool
	sync.Pool
	Err error
}

// Run ...
func (t *task) Run() {
	runs := atomic.NewInt64(0)
	ctx, cancel := context.WithCancel(context.TODO())
	for {
		if t.stop.Load() {
			cancel()
			break
		}
		if runs.Load() <= t.MaxLimit {
			runs.Add(1)
			go func() {
				defer runs.Add(-1)
				for v := t.Get(); v != nil; {
					if fn, b := v.(CallFunc); b {
						fn(ctx)
					}
				}
			}()
		}
	}

	for {
		if runs.Load() == 0 {
			return
		}
		//only wait 3 sec
		time.Sleep(3 * time.Second)
		return
	}

}

// AddCall ...
func (t *task) AddCall(callFunc CallFunc) {
	t.Put(callFunc)
}

// Waiting ...
func (t *task) Waiting() {
	for {
		if t.stop.Load() {
			return
		}
		//every 5 sec check once
		time.Sleep(5 * time.Second)
	}
}

// CallFunc ...
type CallFunc func(ctx context.Context)

// Task ...
type Task interface {
	Run()
	AddCall(callFunc CallFunc)
	Waiting()
}

// New ...
func New() Task {
	return &task{
		stop:     atomic.NewBool(false),
		MaxLimit: 0,
		Pool:     sync.Pool{},
	}
}
