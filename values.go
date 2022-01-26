package main

import (
    `fmt`
    `math`
    `strconv`
    `strings`
)

type Value interface {
    String() string
    IsIdentity() bool
}

type Callable interface {
    Value
    Call([]Value) Value
}

type Numerical interface {
    Value
    Kind() NumKind
    AsInt() Int
    AsFloat() Float
    AsComplex() Complex
}

func AsList(v Value) (*List, bool) {
    r, ok := v.(*List)
    return r, ok || v == nil
}

func AsString(v Value) string {
    if v == nil {
        return "()"
    } else {
        return v.String()
    }
}

func AsDisplay(v Value) string {
    switch vv := v.(type) {
        case nil    : return "()"
        case Char   : return string(vv)
        case String : return string(vv)
        default     : return v.String()
    }
}

type (
    Int     int64
    Bool    bool
    Char    rune
    Atom    string
    Float   float64
    String  string
    Complex complex128
)

type List struct {
    Car Value
    Cdr Value
}

func MakeList(vals ...Value) *List {
    var p, q *List
    for _, v := range vals { AppendValue(&p, &q, v) }
    return p
}

func MakePair(car Value, cdr Value) *List {
    return &List {
        Car: car,
        Cdr: cdr,
    }
}

func AppendValue(p **List, q **List, v Value) {
    if *p == nil {
        *p = new(List)
        *q, (*p).Car = *p, v
    } else {
        r := new(List)
        r.Car, (*q).Cdr, *q = v, r, r
    }
}

/** Value Protocol **/

func (Int)     IsIdentity() bool { return true  }
func (Bool)    IsIdentity() bool { return true  }
func (Char)    IsIdentity() bool { return true  }
func (Atom)    IsIdentity() bool { return false }
func (*List)   IsIdentity() bool { return false }
func (Float)   IsIdentity() bool { return true  }
func (String)  IsIdentity() bool { return true  }
func (Complex) IsIdentity() bool { return true  }

func (self Int) String() string {
    return strconv.Itoa(int(self))
}

func (self Bool) String() string {
    if self {
        return "#t"
    } else {
        return "#f"
    }
}

func (self Char) String() string {
    switch self {
        case ' '  : return `#\space`
        case '\n' : return `#\newline`
        case '\b' : return `#\backspace`
        case '\t' : return `#\tab`
        case '\f' : return `#\page`
        case '\r' : return `#\return`
        case 0x7f : return `#\rubout`
        default   : return `#\` + string(self)
    }
}

func (self Atom) String() string {
    return string(self)
}

func (self *List) String() string {
    var ok bool
    var vv *List
    var rb []string

    /* nil list */
    if vv = self; vv == nil {
        return "()"
    }

    /* dump every element */
    for vv != nil {
        d := vv.Cdr
        s := AsString(vv.Car)
        vv, ok = AsList(vv.Cdr)

        /* also append the last item if this is not a proper list */
        if rb = append(rb, s); !ok {
            rb = append(rb, ".", AsString(d))
        }
    }

    /* compose the final list */
    return fmt.Sprintf(
        "(%s)",
        strings.Join(rb, " "),
    )
}

func (self Float) String() string {
    vv := strconv.FormatFloat(float64(self), 'g', -1, 64)
    vp := strings.Split(vv, "e")

    /* already contains a decimal point */
    if strings.ContainsRune(vp[0], '.') {
        return vv
    }

    /* add one if needed */
    vp[0] += ".0"
    return strings.Join(vp, "e")
}

func (self String) String() string {
    return strconv.Quote(string(self))
}

func (self Complex) String() string {
    if im := imag(complex128(self)); im >= 0 {
        return fmt.Sprintf("%g+%gi", real(complex128(self)), im)
    } else {
        return fmt.Sprintf("%g-%gi", real(complex128(self)), -im)
    }
}

/** Numerical Protocols for Int **/

func (self Int) Kind()      NumKind { return NumInt }
func (self Int) AsInt()     Int     { return self }
func (self Int) AsFloat()   Float   { return Float(self) }
func (self Int) AsComplex() Complex { return Complex(complex(float64(self), 0)) }

/** Numerical Protocols for Float **/

func (self Float) Kind()      NumKind { return NumFloat }
func (self Float) AsInt()     Int     { return Int(self) }
func (self Float) AsFloat()   Float   { return self }
func (self Float) AsComplex() Complex { return Complex(complex(float64(self), 0)) }

/** Numerical Protocols for Complex **/

func (self Complex) Kind()      NumKind { return NumComplex }
func (self Complex) AsInt()     Int     { return Int(self.AsRealNumber()) }
func (self Complex) AsFloat()   Float   { return Float(self.AsRealNumber()) }
func (self Complex) AsComplex() Complex { return self }

func (self Complex) Magnitude() Float {
    return Float(math.Hypot(
        real(complex128(self)),
        imag(complex128(self)),
    ))
}

func (self Complex) AsRealNumber() float64 {
    if v := complex128(self); imag(v) != 0 {
        panic("cast: cannot convert complex numbers with non-zero imaginary part into real numbers")
    } else {
        return real(v)
    }
}
