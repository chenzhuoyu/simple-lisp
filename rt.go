package main

import (
    `unsafe`
)

type _GoIface struct {
    itab unsafe.Pointer
    data unsafe.Pointer
}

func (self _GoIface) pack() Value {
    return *(*Value)(unsafe.Pointer(&self))
}

type _GoString struct {
    ptr unsafe.Pointer
    len int
}

func (self _GoString) String() string {
    return *(*string)(unsafe.Pointer(&self))
}

func mkstr(p unsafe.Pointer, n int) _GoString {
    return _GoString {
        ptr: p,
        len: n,
    }
}

func mkval(t unsafe.Pointer, p unsafe.Pointer) _GoIface {
    return _GoIface {
        itab: t,
        data: p,
    }
}

func valitab(v Value) unsafe.Pointer {
    return (*_GoIface)(unsafe.Pointer(&v)).itab
}

func valaddr(v Value) unsafe.Pointer {
    return (*_GoIface)(unsafe.Pointer(&v)).data
}

func straddr(s string) unsafe.Pointer {
    return (*_GoString)(unsafe.Pointer(&s)).ptr
}
