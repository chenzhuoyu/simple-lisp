package main

import (
    `fmt`
    `math/big`
    `strconv`
    `strings`
)

type Value interface {
    String() string
    IsIdentity() bool
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

type (
    Bool    bool
    Char    rune
    Atom    string
    String  string
    Number  float64
    Complex complex128
)

type Int struct {
    *big.Int
}

type Frac struct {
    *big.Rat
}

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

type Proc struct {
    Name  string
    Code  Program
    Args  []string
    Apply func([]Value) Value
}

func (Int)     IsIdentity() bool { return true }
func (Frac)    IsIdentity() bool { return true }
func (Bool)    IsIdentity() bool { return true }
func (Char)    IsIdentity() bool { return true }
func (Atom)    IsIdentity() bool { return false }
func (*List)   IsIdentity() bool { return false }
func (*Proc)   IsIdentity() bool { return true }
func (String)  IsIdentity() bool { return true }
func (Number)  IsIdentity() bool { return true }
func (Complex) IsIdentity() bool { return true }

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

func (self *Proc) String() string {
    if len(self.Args) == 0 {
        return fmt.Sprintf("#[proc (%s)]", self.Name)
    } else {
        return fmt.Sprintf("#[proc (%s %s)]", self.Name, strings.Join(self.Args, " "))
    }
}

func (self String) String() string {
    return strconv.Quote(string(self))
}

func (self Number) String() string {
    return strconv.FormatFloat(float64(self), 'g', -1, 64)
}

func (self Complex) String() string {
    return fmt.Sprintf("%g+%gi", real(complex128(self)), imag(complex128(self)))
}
