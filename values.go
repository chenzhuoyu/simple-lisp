package main

import (
    `fmt`
    `math`
    `math/big`
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
    Bool    bool
    Char    rune
    Atom    string
    String  string
    Double  float64
    Complex complex128
)

func (self Double) Exact() Value {
    if vv := new(big.Rat).SetFloat64(float64(self)); vv.IsInt() {
        return Int{vv.Num()}
    } else {
        return Frac{vv}
    }
}

func (self Double) Round() Double {
    return Double(math.RoundToEven(float64(self)))
}

func (self Complex) Exact() Value {
    if x := complex128(self); imag(x) != 0 {
        panic("exact: cannot measure the exact value of complex numbers with non-zero imaginary part")
    } else {
        return Double(real(x)).Exact()
    }
}

func (self Complex) Round() Double {
    if x := complex128(self); imag(x) != 0 {
        panic("round: cannot round complex numbers with non-zero imaginary part")
    } else {
        return Double(real(x)).Round()
    }
}

type Int struct {
    *big.Int
}

func MakeInt(v int64) Int {
    return Int {
        new(big.Int).SetInt64(v),
    }
}

func (self Int) Cmp(v Int) int {
    return self.Int.Cmp(v.Int)
}

func (self Int) Add(v Int) Int {
    return Int {
        new(big.Int).Add(self.Int, v.Int),
    }
}

func (self Int) Sub(v Int) Int {
    return Int {
        new(big.Int).Sub(self.Int, v.Int),
    }
}

func (self Int) Mul(v Int) Int {
    return Int {
        new(big.Int).Mul(self.Int, v.Int),
    }
}

func (self Int) Div(v Int) Int {
    return Int {
        new(big.Int).Quo(self.Int, v.Int),
    }
}

func (self Int) Mod(v Int) Int {
    return Int {
        new(big.Int).Rem(self.Int, v.Int),
    }
}

func (self Int) Frac() Frac {
    return Frac {
        new(big.Rat).SetInt(self.Int),
    }
}

func (self Int) Double() Double {
    v, _ := new(big.Float).SetInt(self.Int).Float64()
    return Double(v)
}

func (self Int) Complex() Complex {
    v, _ := new(big.Float).SetInt(self.Int).Float64()
    return Complex(complex(v, 0))
}

type Frac struct {
    *big.Rat
}

func MakeFrac(a Int, b Int) Frac {
    return Frac {
        new(big.Rat).SetFrac(a.Int, b.Int),
    }
}

func (self Frac) Cmp(v Frac) int {
    return self.Rat.Cmp(v.Rat)
}

func (self Frac) Add(v Frac) Frac {
    return Frac {
        new(big.Rat).Add(self.Rat, v.Rat),
    }
}

func (self Frac) Sub(v Frac) Frac {
    return Frac {
        new(big.Rat).Sub(self.Rat, v.Rat),
    }
}

func (self Frac) Mul(v Frac) Frac {
    return Frac {
        new(big.Rat).Mul(self.Rat, v.Rat),
    }
}

func (self Frac) Div(v Frac) Frac {
    return Frac {
        new(big.Rat).Quo(self.Rat, v.Rat),
    }
}

func (self Frac) Round() Int {
    return Int {
        new(big.Int).Quo(self.Num(), self.Denom()),
    }
}

func (self Frac) Double() Double {
    v, _ := self.Float64()
    return Double(v)
}

func (self Frac) Complex() Complex {
    v, _ := self.Float64()
    return Complex(complex(v, 0))
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

func (Int)     IsIdentity() bool { return true }
func (Frac)    IsIdentity() bool { return true }
func (Bool)    IsIdentity() bool { return true }
func (Char)    IsIdentity() bool { return true }
func (Atom)    IsIdentity() bool { return false }
func (*List)   IsIdentity() bool { return false }
func (String)  IsIdentity() bool { return true }
func (Double)  IsIdentity() bool { return true }
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

func (self String) String() string {
    return strconv.Quote(string(self))
}

func (self Double) String() string {
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

func (self Complex) String() string {
    return fmt.Sprintf("%g+%gi", real(complex128(self)), imag(complex128(self)))
}
