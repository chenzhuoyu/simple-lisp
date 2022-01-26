package main

import (
    `math`
)

type NumKind uint8

const (
    NumInt NumKind = iota
    NumFloat
    NumComplex
)

func (self NumKind) Coerce(vt NumKind) NumKind {
    if vt > self {
        return vt
    } else {
        return self
    }
}

/** Number Converting **/

func AsNumber(v Value) Numerical {
    if r, ok := v.(Numerical); !ok {
        panic("eval: object is not a number: " + AsString(v))
    } else {
        return r
    }
}

func AsNumbers(v1 Value, v2 Value) (Numerical, Numerical, NumKind) {
    r1, r2 := AsNumber(v1), AsNumber(v2)
    return r1, r2, r1.Kind().Coerce(r2.Kind())
}

/** Number Arithmetic **/

func NumberNeg(v Value) Value {
    switch x := AsNumber(v); x.Kind() {
        case NumInt     : return -x.AsInt()
        case NumFloat   : return -x.AsFloat()
        case NumComplex : return -x.AsComplex()
        default         : panic("-: unreachable")
    }
}

func NumberInv(v Value) Value {
    switch x := AsNumber(v); x.Kind() {
        case NumInt     : fallthrough
        case NumFloat   : return 1.0 / x.AsFloat()
        case NumComplex : return 1.0 / x.AsComplex()
        default         : panic("-: unreachable")
    }
}

func NumberAdd(a Value, b Value) Value {
    switch x, y, vt := AsNumbers(a, b); vt {
        case NumInt     : return x.AsInt() + y.AsInt()
        case NumFloat   : return x.AsFloat() + y.AsFloat()
        case NumComplex : return x.AsComplex() + y.AsComplex()
        default         : panic("+: unreachable")
    }
}

func NumberSub(a Value, b Value) Value {
    switch x, y, vt := AsNumbers(a, b); vt {
        case NumInt     : return x.AsInt() - y.AsInt()
        case NumFloat   : return x.AsFloat() - y.AsFloat()
        case NumComplex : return x.AsComplex() - y.AsComplex()
        default         : panic("-: unreachable")
    }
}

func NumberMul(a Value, b Value) Value {
    switch x, y, vt := AsNumbers(a, b); vt {
        case NumInt     : return x.AsInt() * y.AsInt()
        case NumFloat   : return x.AsFloat() * y.AsFloat()
        case NumComplex : return x.AsComplex() * y.AsComplex()
        default         : panic("*: unreachable")
    }
}

func NumberDiv(a Value, b Value) Value {
    switch x, y, vt := AsNumbers(a, b); vt {
        case NumInt     : fallthrough
        case NumFloat   : return x.AsFloat() / y.AsFloat()
        case NumComplex : return x.AsComplex() / y.AsComplex()
        default         : panic("/: unreachable")
    }
}

func NumberRound(v Value) Value {
    switch x := AsNumber(v); x.Kind() {
        case NumInt     : return v
        case NumFloat   : fallthrough
        case NumComplex : return Float(math.RoundToEven(float64(x.AsFloat())))
        default         : panic("inexact->exact: unreachable")
    }
}

func NumberMagnitude(v Value) Value {
    switch x := AsNumber(v); x.Kind() {
        case NumInt     : fallthrough
        case NumFloat   : return v
        case NumComplex : return x.AsComplex().Magnitude()
        default         : panic("inexact->exact: unreachable")
    }
}

func NumberCompareEq(a Value, b Value) bool {
    switch x, y, vt := AsNumbers(a, b); vt {
        case NumInt     : return x.AsInt() == y.AsInt()
        case NumFloat   : return x.AsFloat() == y.AsFloat()
        case NumComplex : return x.AsComplex() == y.AsComplex()
        default         : panic("=: unreachable")
    }
}

func NumberCompareLt(a Value, b Value) bool {
    switch x, y, vt := AsNumbers(a, b); vt {
        case NumInt     : return x.AsInt() < y.AsInt()
        case NumFloat   : return x.AsFloat() < y.AsFloat()
        case NumComplex : panic("<: complex numbers can only be compared for equality")
        default         : panic("<: unreachable")
    }
}

func NumberCompareGt(a Value, b Value) bool {
    switch x, y, vt := AsNumbers(a, b); vt {
        case NumInt     : return x.AsInt() > y.AsInt()
        case NumFloat   : return x.AsFloat() > y.AsFloat()
        case NumComplex : panic(">: complex numbers can only be compared for equality")
        default         : panic(">: unreachable")
    }
}

func NumberCompareLte(a Value, b Value) bool {
    switch x, y, vt := AsNumbers(a, b); vt {
        case NumInt     : return x.AsInt() <= y.AsInt()
        case NumFloat   : return x.AsFloat() <= y.AsFloat()
        case NumComplex : panic("<=: complex numbers can only be compared for equality")
        default         : panic("<=: unreachable")
    }
}

func NumberCompareGte(a Value, b Value) bool {
    switch x, y, vt := AsNumbers(a, b); vt {
        case NumInt     : return x.AsInt() >= y.AsInt()
        case NumFloat   : return x.AsFloat() >= y.AsFloat()
        case NumComplex : panic(">=: complex numbers can only be compared for equality")
        default         : panic(">=: unreachable")
    }
}
