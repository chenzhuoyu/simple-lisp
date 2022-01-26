package main

import (
    `os`
)

const (
    MaxBufferSize = 65536
)

type BufferedFileWriter struct {
    wb []byte
    fp *os.File
}

func (self *BufferedFileWriter) Write(p []byte) (int, error) {
    var ret int
    var rem int

    /* write each part */
    for len(p) > 0 {
        nb := len(p)
        wb := len(self.wb)

        /* no space left in buffer, attempt to flush the buffer */
        if wb == MaxBufferSize {
            if err := self.Flush(); err != nil {
                return ret, err
            }
        }

        /* remaining buffer space */
        wb = len(self.wb)
        rem = MaxBufferSize - wb

        /* check if the entire chunk can fit into the output buffer */
        if nb < rem {
            rem = nb
        }

        /* write part of the data into output buffer */
        ret += rem
        self.wb, p = append(self.wb, p[:rem]...), p[rem:]
    }

    /* all done */
    return ret, nil
}

func (self *BufferedFileWriter) Close() error {
    if err := self.Flush(); err != nil {
        return err
    } else {
        return self.fp.Close()
    }
}

func (self *BufferedFileWriter) Flush() error {
    var out int
    var err error

    /* write all pending data into file */
    if out, err = self.fp.Write(self.wb); err != nil {
        return err
    }

    /* move remaining data to front */
    copy(self.wb[:len(self.wb) - out], self.wb[out:])
    self.wb = self.wb[:len(self.wb) - out:cap(self.wb)]
    return nil
}

func CreateBufferedWriter(fp *os.File) *BufferedFileWriter {
    return &BufferedFileWriter {
        fp: fp,
        wb: make([]byte, 0, MaxBufferSize),
    }
}
