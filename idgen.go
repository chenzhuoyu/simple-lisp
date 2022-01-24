package main

import (
    `sync/atomic`
)

var (
    incr uint32
)

func nextid() int {
    return int(atomic.AddUint32(&incr, 1))
}
