package main

import (
    `fmt`
    `os`
)

type Port struct {
    name string
    file FileLike
}

type FileLike interface {
    Write(p []byte) (int, error)
    Close() error
}

var PortStdout = &Port {
    name: "<stdout>",
    file: os.Stdout,
}

func CreatePort(name string, file FileLike) *Port {
    return &Port {
        name: name,
        file: file,
    }
}

func OpenFileWritePort(fname string) *Port {
    if fp, err := os.OpenFile(fname, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666); err != nil {
        panic(fmt.Sprintf("port: cannot open %s for write: %s", fname, err))
    } else {
        return CreatePort(fname, CreateBufferedWriter(fp))
    }
}

func (self *Port) Close() {
    _ = self.file.Close()
}

func (self *Port) Write(v []byte) {
    if _, err := self.file.Write(v); err != nil {
        panic(fmt.Sprintf("port: write error to port %s: %s", self.name, err))
    }
}

func (self *Port) String() string {
    return fmt.Sprintf("#[port-%s]", self.name)
}

func (self *Port) IsIdentity() bool {
    return true
}
