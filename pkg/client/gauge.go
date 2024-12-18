package client

import (
	"sync/atomic"
	"time"
)

type Gauge struct {
	name  string
	value int64
}

func NewGauge(name string) *Gauge {
	return &Gauge{
		name: name,
	}
}

func (g *Gauge) Set(value int64) {
	atomic.SwapInt64(&g.value, value)
}

func (c *Gauge) getValue() instance {
	return instance{
		Name:        c.name,
		Value:       float64(c.value),
		MessageType: "GAUGE",
		Time:        time.Now().UTC().UnixMilli(),
	}
}
