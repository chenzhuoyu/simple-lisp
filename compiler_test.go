package main

import (
    `testing`
)

func stmt(s string) *List {
    return CreateParser(s).Parse().Cdr.(*List).Car.(*List)
}

func TestCompiler_Desugar(t *testing.T) {
    tests := []struct{
        fn func(Compiler, *List) *List
        ss string
    } {{
        fn: Compiler.desugarDo,
        ss: `(do ((a 0 (+ a 1))
                  (b 0 (+ b 1)))
                 ((and (> a 10) (> b 10)) (+ a b))
                (display a)
                (display #\space)
                (display b)
                (newline))`,
    }, {
        fn: func(c Compiler, s *List) *List { return c.desugarLet(s, Let) },
        ss: `(let ((a 1) (b 2)) (display (+ a b)) (newline))`,
    }, {
        fn: func(c Compiler, s *List) *List { return c.desugarLet(s, LetRec) },
        ss: `(letrec ((fac (lambda (v r)
                (if (= v 0) r (fac (- v 1) (* v r))))))
                (display (fac 10 1))
                (newline))`,
    }, {
        fn: func(c Compiler, s *List) *List { return c.desugarLet(s, LetStar) },
        ss: `(let ((a 1) (b 2)) (display (+ a b)) (newline))`,
    }}
    for _, ts := range tests {
        println(ts.fn(Compiler{}, stmt(ts.ss).Cdr.(*List)).String())
    }
}

func TestCompiler_Compile(t *testing.T) {
    src := `
        (let ((r 1))
            (do ((i 1 (+ i 1))) ((> i 10) r)
                (set! r (* r i)))
            (display r)
            (newline))
    `
    println(Compiler{}.Compile(CreateParser(src).Parse()).String())
}
