package main

import (
    `testing`
)

func stmt(s string) *List {
    return CreateParser(s).Parse().Cdr.(*List).Car.(*List)
}

func TestCompiler_Desugar(t *testing.T) {
    println(Compiler{}.desugarDo(stmt(`(do ((i 0 (+ i 1))) ((> i 10)) (display i) (newline))`).Cdr.(*List)).String())
    println(Compiler{}.desugarLet(stmt(`(let ((a 1) (b 2)) (display (+ a b)) (newline))`).Cdr.(*List), true).String())
    println(Compiler{}.desugarLet(stmt(`(let ((a 1) (b 2)) (display (+ a b)) (newline))`).Cdr.(*List), false).String())
    println(Compiler{}.desugarBegin(stmt(`(begin (display 123) (newline))`).Cdr.(*List)).String())
}

func TestCompiler_Compile(t *testing.T) {
    println(Compiler{}.Compile(CreateParser(`(+ 1 (car (cdr (quote (cons 1 2)))))`).Parse()).String())
}
