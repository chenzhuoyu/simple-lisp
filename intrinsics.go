package main

import (
    `fmt`
)

type Intrinsic struct {
    Name string
    Proc func([]Value) Value
}

var (
	intrinsicsTab = make(map[string]*Intrinsic)
)

func newIntrinsic(name string, proc func([]Value) Value) *Intrinsic {
    return &Intrinsic {
        Name: name,
        Proc: proc,
    }
}

func RegisterIntrinsic(name string, proc func([]Value) Value) {
    if _, ok := intrinsicsTab[name]; ok {
        panic("registry: duplicated intrinsic proc: " + name)
    } else {
        intrinsicsTab[name] = newIntrinsic(name, proc)
    }
}

func (self *Intrinsic) Call(args []Value) Value {
    return self.Proc(args)
}

func (self *Intrinsic) String() string {
    return fmt.Sprintf("#[intrinsic-%s]", self.Name)
}

func (self *Intrinsic) IsIdentity() bool {
    return true
}

/** Arithmetic Operators **/

func intrinsicsAdd(args []Value) Value {
    switch len(args) {
        case 0  : return Int(0)
        case 1  : return AsNumber(args[0])
        case 2  : return NumberAdd(args[0], args[1])
        default : return reduceSequential(args, NumberAdd)
    }
}

func intrinsicsSub(args []Value) Value {
    switch len(args) {
        case 0  : panic("-: proc requies at least 1 argument")
        case 1  : return NumberNeg(args[0])
        case 2  : return NumberSub(args[0], args[1])
        default : return reduceSequential(args, NumberSub)
    }
}

func intrinsicsMul(args []Value) Value {
    switch len(args) {
        case 0  : return Int(1)
        case 1  : return AsNumber(args[0])
        case 2  : return NumberMul(args[0], args[1])
        default : return reduceSequential(args, NumberMul)
    }
}

func intrinsicsDiv(args []Value) Value {
    switch len(args) {
        case 0  : panic("/: proc requies at least 1 argument")
        case 1  : return NumberInv(args[0])
        case 2  : return NumberDiv(args[0], args[1])
        default : return reduceSequential(args, NumberDiv)
    }
}

func init() {
    RegisterIntrinsic("+", intrinsicsAdd)
    RegisterIntrinsic("-", intrinsicsSub)
    RegisterIntrinsic("*", intrinsicsMul)
    RegisterIntrinsic("/", intrinsicsDiv)
}

/** Comparison Operators **/

func intrinsicsEq(args []Value) Value {
    switch len(args) {
        case 0  : fallthrough
        case 1  : return Bool(true)
        case 2  : return Bool(NumberCompareEq(args[0], args[1]))
        default : return Bool(reduceConvolution(args, NumberCompareEq))
    }
}

func intrinsicsLt(args []Value) Value {
    switch len(args) {
        case 0  : fallthrough
        case 1  : return Bool(true)
        case 2  : return Bool(NumberCompareLt(args[0], args[1]))
        default : return Bool(reduceConvolution(args, NumberCompareLt))
    }
}

func intrinsicsGt(args []Value) Value {
    switch len(args) {
        case 0  : fallthrough
        case 1  : return Bool(true)
        case 2  : return Bool(NumberCompareGt(args[0], args[1]))
        default : return Bool(reduceConvolution(args, NumberCompareGt))
    }
}

func intrinsicsLte(args []Value) Value {
    switch len(args) {
        case 0  : fallthrough
        case 1  : return Bool(true)
        case 2  : return Bool(NumberCompareLte(args[0], args[1]))
        default : return Bool(reduceConvolution(args, NumberCompareLte))
    }
}

func intrinsicsGte(args []Value) Value {
    switch len(args) {
        case 0  : fallthrough
        case 1  : return Bool(true)
        case 2  : return Bool(NumberCompareGte(args[0], args[1]))
        default : return Bool(reduceConvolution(args, NumberCompareGte))
    }
}

func init() {
    RegisterIntrinsic("=", intrinsicsEq)
    RegisterIntrinsic("<", intrinsicsLt)
    RegisterIntrinsic(">", intrinsicsGt)
    RegisterIntrinsic("<=", intrinsicsLte)
    RegisterIntrinsic(">=", intrinsicsGte)
}

/** Unary Arithmetic Functions **/

func intrinsicsRound(args []Value) Value {
    if len(args) != 1 {
        panic("round: proc takes exact 1 argument")
    } else {
        return NumberRound(args[0])
    }
}

func intrinsicsMagnitude(args []Value) Value {
    if len(args) != 1 {
        panic("magnitude: proc takes exact 1 argument")
    } else {
        return NumberMagnitude(args[0])
    }
}

func intrinsicsInexactToExact(args []Value) Value {
    if len(args) != 1 {
        panic("inexact->exact: proc takes exact 1 argument")
    } else {
        return AsNumber(args[0]).AsInt()
    }
}

func init() {
    RegisterIntrinsic("round", intrinsicsRound)
    RegisterIntrinsic("magnitude", intrinsicsMagnitude)
    RegisterIntrinsic("inexact->exact", intrinsicsInexactToExact)
}

/** Binary Arithmetic Functions **/

func intrinsicsModulo(args []Value) Value {
    if len(args) != 2 {
        panic("modulo: proc takes exact 2 arguments")
    } else if va, ok := args[0].(Int); !ok {
        panic("modulo: object is not an integer: " + AsString(args[0]))
    } else if vb, ok := args[1].(Int); !ok {
        panic("modulo: object is not an integer: " + AsString(args[1]))
    } else {
        return va % vb
    }
}

func intrinsicsQuotient(args []Value) Value {
    if len(args) != 2 {
        panic("quotient: proc takes exact 2 arguments")
    } else if va, ok := args[0].(Int); !ok {
        panic("quotient: object is not an integer: " + AsString(args[0]))
    } else if vb, ok := args[1].(Int); !ok {
        panic("quotient: object is not an integer: " + AsString(args[1]))
    } else {
        return va / vb
    }
}

func init() {
    RegisterIntrinsic("modulo", intrinsicsModulo)
    RegisterIntrinsic("quotient", intrinsicsQuotient)
}

/** Value Constructors **/

func intrinsicMakeRectangular(args []Value) Value {
    if len(args) != 2 {
        panic("make-rectangular: proc takes exact 2 arguments")
    } else {
        return Complex(complex(float64(AsNumber(args[0]).AsFloat()), float64(AsNumber(args[1]).AsFloat())))
    }
}

func init() {
    RegisterIntrinsic("make-rectangular", intrinsicMakeRectangular)
}

/** Input / Output Functions **/

func intrinsicsDisplay(args []Value) Value {
    var ok bool
    var wp *Port

    /* check for arguments */
    if len(args) != 1 && len(args) != 2 {
        panic("display: proc requires 1 or 2 arguments")
    }

    /* check for optional port */
    if wp = PortStdout; len(args) == 2 {
        if wp, ok = args[1].(*Port); !ok {
            panic("display: object is not a port: " + AsString(args[1]))
        }
    }

    /* display the value */
    wp.Write([]byte(AsDisplay(args[0])))
    return nil
}

func intrinsicsNewline(args []Value) Value {
    var ok bool
    var wp *Port

    /* check for arguments */
    if len(args) != 0 && len(args) != 1 {
        panic("newline: proc requires 1 or 2 arguments")
    }

    /* check for optional port */
    if wp = PortStdout; len(args) == 1 {
        if wp, ok = args[0].(*Port); !ok {
            panic("newline: object is not a port: " + AsString(args[0]))
        }
    }

    /* display the newline */
    wp.Write([]byte{'\n'})
    return nil
}

func intrinsicsCallWithOutputFile(args []Value) Value {
    var ok bool
    var fn String
    var cb LoadedProc

    /* extract the file name and callback */
    if len(args) != 2                     { panic("call-with-output-file: proc requires exact 2 arguments") }
    if fn, ok = args[0].(String)    ; !ok { panic("call-with-output-file: object is not a string: " + AsString(args[0])) }
    if cb, ok = args[1].(LoadedProc); !ok { panic("call-with-output-file: object is not a callable proc: " + AsString(args[1])) }

    /* open a new port */
    file := string(fn)
    port := OpenFileWritePort(file)

    /* call the function with the port */
    defer port.Close()
    return cb.Call([]Value{port})
}

func init() {
    RegisterIntrinsic("display", intrinsicsDisplay)
    RegisterIntrinsic("newline", intrinsicsNewline)
    RegisterIntrinsic("call-with-output-file", intrinsicsCallWithOutputFile)
}
