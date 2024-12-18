package client

import (
	"sync/atomic"
	"time"
)

type Counter struct {
	name  string
	value int64
}

func NewCounter(name string) *Counter {
	return &Counter{
		name: name,
	}
}

func (c *Counter) Inc() {
	atomic.AddInt64(&c.value, 1)
}

func (c *Counter) getValue() instance {
	return instance{
		Name:        c.name,
		Value:       float64(c.value),
		MessageType: "COUNTER",
		Time:        time.Now().UTC().UnixMilli(),
	}
}
