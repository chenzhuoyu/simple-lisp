package main

import (
    `fmt`
    `os`
)

func main() {
    var err error
    var src []byte
    var ins Compiler

    /* check for args */
    if len(os.Args) != 2 || os.Args[1] == "-h" {
        println(fmt.Sprintf("usage: %s [-h] [file-name]", os.Args[0]))
        os.Exit(1)
    }

    /* read the source */
    if src, err = os.ReadFile(os.Args[1]); err != nil {
        panic(fmt.Sprintf("slisp: unable to read %s: %s", os.Args[1], err))
    }

    /* evaluate the source */
    Evaluate(
        CreateGlobalScope(),
        ins.Compile(CreateParser(string(src)).Parse()),
    )
}
