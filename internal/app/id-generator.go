package app

import (
	"sync/atomic"
	"time"

	"github.com/MaratBR/openlibrary/internal/commonutil"
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

func GenID() int64 {
	return defaultIdGenerator.get()
}

func genOpaqueID() string {
	s, err := commonutil.GenerateRandomStringURLSafe(32)
	if err != nil {
		panic(err)
	}
	return s
}
