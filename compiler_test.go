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
        ss: `(letrec ((fac (λ (v r)
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
    // src := `
    //     (let ((r 1))
    //         (do ((i 1 (+ i 1))) ((> i 10) r)
    //             (set! r (* r i)))
    //         (display r)
    //         (newline))
    // `
    src := `
(define real-min -2.5)
(define real-max +1.5)
(define imag-min -2.0)
(define imag-max +2.0)

(define max-color 255)
(define iter-count 40)
(define escape-radius 4.0)

(define width 1024)
(define scale (/ width (- real-max real-min)))
(define height (inexact->exact (round (* scale (- imag-max imag-min)))))

(define (complex real imag)
  (make-rectangular real imag))
`
    println(Compiler{}.Compile(CreateParser(src).Parse()).String())
}
