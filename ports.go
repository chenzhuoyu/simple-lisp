package main

import (
    `fmt`
    `io`
    `os`
)

type Port struct {
    name string
    file io.ReadWriteCloser
}

var PortStdout = &Port {
    name: "<stdout>",
    file: os.Stdout,
}

func CreatePort(name string, file io.ReadWriteCloser) *Port {
    return &Port {
        name: name,
        file: file,
    }
}

func OpenFileWritePort(fname string) *Port {
    if fp, err := os.OpenFile(fname, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666); err != nil {
        panic(fmt.Sprintf("port: cannot open %s for write: %s", fname, err))
    } else {
        return CreatePort(fname, AsBuffered(fp))
    }
}

func (self *Port) Close() {
    if err := self.file.Close(); err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "warn: cannot close %s, ignored.", self)
    }
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
