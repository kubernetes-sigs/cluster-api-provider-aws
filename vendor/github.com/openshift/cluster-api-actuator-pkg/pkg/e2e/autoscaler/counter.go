package autoscaler

import (
	"sync/atomic"
)

type u32counter uint32

func (c *u32counter) increment() uint32 {
	return atomic.AddUint32((*uint32)(c), 1)
}

func (c *u32counter) decrement() uint32 {
	return atomic.AddUint32((*uint32)(c), ^uint32(0))
}

func (c *u32counter) get() uint32 {
	return atomic.LoadUint32((*uint32)(c))
}
