package main

import (
    `bufio`
    `os`
)

type BufferedFile struct {
    fp  *os.File
    buf *bufio.ReadWriter
}

func (self BufferedFile) Read(p []byte) (int, error) {
    return self.buf.Read(p)
}

func (self BufferedFile) Write(p []byte) (int, error) {
    return self.buf.Write(p)
}

func (self BufferedFile) Close() error {
    if err := self.buf.Flush(); err != nil {
        return err
    } else {
        return self.fp.Close()
    }
}

func AsBuffered(fp *os.File) BufferedFile {
    return BufferedFile {
        fp  : fp,
        buf : bufio.NewReadWriter(bufio.NewReader(fp), bufio.NewWriter(fp)),
    }
}
