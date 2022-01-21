package main

import (
	"testing"
)

func TestCompiler_Desugar(t *testing.T) {
	p := CreateParser(`(let ((a 1) (b 2)) (begin (display (+ a b)) (newline)))`).Parse()
	println(Compiler{}.desugarLet(p.Cdr.(*List).Car.(*List).Cdr.(*List), true).String())
	println(Compiler{}.desugarLet(p.Cdr.(*List).Car.(*List).Cdr.(*List), false).String())
}

func TestCompiler_Compile(t *testing.T) {
	println(Compiler{}.Compile(CreateParser(`(+ 1 (car (cdr (quote (cons 1 2)))))`).Parse()).String())
}
