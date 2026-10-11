package idgen

import (
	"sync/atomic"
	"time"
)

type idGenerator struct {
	counter atomic.Int32
}

func (g *idGenerator) get() int64 {
	c := uint64(g.counter.Add(1))
	ts := uint64(time.Now().UnixMilli())
	v := ts<<32 | c
	v &= 0b01111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111
	return int64(v)
}

var (
	defaultIdGenerator = new(idGenerator)
)

// GenID uses the shared process-wide counter for application entity IDs.
func GenID() int64 {
	return defaultIdGenerator.get()
}
