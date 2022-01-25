package main

import (
    `math/big`
)

type _Num uint8

const (
    _T_int _Num = iota
    _T_frac
    _T_float
    _T_complex
)

func maxtype(a _Num, b _Num) _Num {
    if a > b {
        return a
    } else {
        return b
    }
}

func numtypeof(v Value) _Num {
    switch v.(type) {
        case Int     : return _T_int
        case Frac    : return _T_frac
        case Double  : return _T_float
        case Complex : return _T_complex
        default      : panic("cast: object is not a number: " + AsString(v))
    }
}

func numasfrac(v Value) Frac {
    switch vv := v.(type) {
        case Int     : return vv.Frac()
        case Frac    : return vv
        case Double  : panic("cast: cannot convert float to frac: " + v.String())
        case Complex : panic("cast: cannot convert complex to frac: " + v.String())
        default      : panic("cast: object is not a number: " + AsString(v))
    }
}

func numasfloat(v Value) Double {
    switch vv := v.(type) {
        case Int     : return vv.Double()
        case Frac    : return vv.Double()
        case Double  : return vv
        case Complex : panic("cast: cannot convert complex to float: " + v.String())
        default      : panic("cast: object is not a number: " + AsString(v))
    }
}

func numascomplex(v Value) Complex {
    switch vv := v.(type) {
        case Int     : return vv.Complex()
        case Frac    : return vv.Complex()
        case Double  : return Complex(complex(float64(vv), 0))
        case Complex : return vv
        default      : panic("cast: object is not a number: " + AsString(v))
    }
}

func inttodouble(v *big.Int) (r float64) {
    f := big.Float{}
    r, _ = f.SetInt(v).Float64()
    return
}

func rattodouble(v *big.Rat) (r float64) {
    r, _ = v.Float64()
    return
}
