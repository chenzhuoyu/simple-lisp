package main

import (
    `testing`
)

func TestEval_Expression(t *testing.T) {
    src := `
        (letrec ((fac (λ (v r)
                         (if (= v 0) r (fac (- v 1) (* v r))))))
                (display (fac 100 1))
                (newline))
    `
    prog := Compiler{}.Compile(CreateParser(src).Parse())
    println(prog.String())
    println(AsString(Evaluate(CreateGlobalScope(), prog)))
}
