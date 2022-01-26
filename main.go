package main

import (
    `fmt`
    `os`
)

func readfile(fname string) string {
    var nbr int
    var err error
    var rfp *os.File

    /* open the file */
    if rfp, err = os.OpenFile(fname, os.O_RDONLY, 0); err != nil {
        panic(fmt.Sprintf("io: unable to open %s: %s", fname, err))
    }

    /* allocate memory for file */
    buf := make([]byte, MaxBufferSize)
    ret := make([]byte, 0, MaxBufferSize)

    /* read all data */
    for err == nil {
        if nbr, err = rfp.Read(buf); nbr != 0 {
            ret = append(ret, buf[:nbr]...)
        }
    }

    /* close the file, and convert the result to string */
    _ = rfp.Close()
    return string(ret)
}

func main() {
    if len(os.Args) != 2 || os.Args[1] == "-h" {
        println(fmt.Sprintf("usage: %s [-h] [file-name]", os.Args[0]))
    } else {
        Evaluate(CreateGlobalScope(), Compiler{}.Compile(CreateParser(readfile(os.Args[1])).Parse()))
    }
}
