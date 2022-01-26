package main

var (
    incr uint32
)

func nextid() int {
    incr++
    return int(incr)
}
