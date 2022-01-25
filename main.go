package main

import (
    `bufio`
    `fmt`
    `os`
    `strings`
)

func repl() {
    println("Simple Lisp Interpreter")
    println()

    /* create a global scope */
    ctx := CreateGlobalScope()
    src := bufio.NewReader(os.Stdin)

    /* read & interpret every line */
    for {
        if _, err := os.Stdout.WriteString("> "); err != nil {
            break
        } else if line, _, err := src.ReadLine(); err != nil {
            break
        } else if rbuf := strings.TrimSpace(string(line)); rbuf == ",q" {
            break
        } else if rbuf != "" {
            println(AsString(Evaluate(ctx, Compiler{}.Compile(CreateParser(rbuf).Parse()))))
        }
    }
}

func usage() {
    println(fmt.Sprintf("usage: %s [-h] [file-name]", os.Args[0]))
}

func runScript(fname string) {
    if src, err := os.ReadFile(fname); err != nil {
        panic(fmt.Sprintf("slisp: unable to read %s: %s", fname, err))
    } else {
        Evaluate(CreateGlobalScope(), Compiler{}.Compile(CreateParser(string(src)).Parse()))
    }
}

func main() {
    if len(os.Args) == 1 {
        repl()
    } else if len(os.Args) > 2 || os.Args[1] == "-h" {
        usage()
    } else {
        runScript(os.Args[1])
    }
}
