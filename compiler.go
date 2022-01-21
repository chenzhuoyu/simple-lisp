package main

import (
    `fmt`
    `strconv`
    `strings`
)

type (
	OpCode   uint8
    Program  []Instr
    Compiler struct{}
)

const (
    Lambda = "λ"
)

const (
    _ OpCode = iota
    OP_ldconst          // ldconst      <val>       : Push <val> onto stack.
    OP_ldvar            // ldvar        <name>      : Push the content of variable <name> onto stack.
    OP_store            // store        <name>      : Store the content of stack top into variable <name>.
    OP_define           // define       <name>      : Define a global variable <name> with content at the stack top.
    OP_car              // car                      : Get the first part of a pair.
    OP_cdr              // cdr                      : Get the second part of a pair.
    OP_cons             // cons                     : Construct a new pair from stack top.
    OP_eval             // eval                     : Evaluate the value on stack top.
    OP_drop             // drop                     : Drop one value from stack.
    OP_enter            // enter                    : Enter a new scope.
    OP_leave            // leave                    : Leave to the parent scope.
    OP_goto             // goto         <pc>        : Goto <pc> unconditionally.
    OP_if_false         // if_false     <pc>        : If the stack top is #f, goto <pc>.
    OP_apply            // apply        <argc>      : Apply procedure on stack with <argc> arguments.
    OP_return           // return                   : Return from procedure.
)

type Instr struct {
    Op OpCode
    Iv uint32
    Rv Value
}

func (self Instr) String() string {
    switch self.Op {
        case OP_ldconst : return fmt.Sprintf("ldconst     %s", self.Rv)
        case OP_ldvar   : return fmt.Sprintf("ldvar       %s", self.Rv)
        case OP_store   : return fmt.Sprintf("store       %s", self.Rv)
        case OP_define  : return fmt.Sprintf("define      %s", self.Rv)
        case OP_car     : return "car"
        case OP_cdr     : return "cdr"
        case OP_cons    : return "cons"
        case OP_eval    : return "eval"
        case OP_drop    : return "drop"
        case OP_enter   : return "enter"
        case OP_leave   : return "leave"
        case OP_goto    : return fmt.Sprintf("goto       @%d", self.Iv)
        case OP_if_false: return fmt.Sprintf("if_false   @%d", self.Iv)
        case OP_apply   : return fmt.Sprintf("apply       %d", self.Iv)
        case OP_return  : return "return"
        default         : return fmt.Sprintf("OpCode(%d)", self)
    }
}

func (self *Program) add(op OpCode)             { *self = append(*self, Instr{Op: op}) }
func (self *Program) val(op OpCode, val Value)  { *self = append(*self, Instr{Op: op, Rv: val}) }
func (self *Program) i32(op OpCode, val uint32) { *self = append(*self, Instr{Op: op, Iv: val}) }

func (self Program) String() string {
    var idx int
    var val Instr
    var buf []string

    /* empty program */
    if len(self) == 0 {
        return "(empty program)"
    }

    /* count the maximum digits */
    nb := len(self) - 1
    nd := len(strconv.Itoa(nb))
    fs := fmt.Sprintf("%%%dd :  %%s\n", nd)

    /* disassemble every instruction */
    for idx, val = range self {
        buf = append(buf, fmt.Sprintf(fs, idx, val))
    }

    /* pack into result */
    return strings.Join(buf, "")
}

func (self Compiler) Compile(src *List) (p Program) {
    self.compileList(&p, src)
    return
}

func (self Compiler) compileList(p *Program, v *List) {
    var ok bool
    var at Atom
    var vv *List

    /* empty list */
    if v == nil {
        p.val(OP_ldconst, nil)
        return
    }

    /* (car v) is not an atom, apply the list immediately */
    if at, ok = v.Car.(Atom); !ok {
        p.i32(OP_apply, self.compileArgs(p, v, -1))
        return
    }

    /* must be a proper list to be applicable */
    if vv, ok = AsList(v.Cdr); !ok {
        panic("compile: improper list is not applicable: " + v.String())
    }

    /* check for built-in atoms */
    switch at {
        case "car"    : self.compileArgs(p, vv, 1); p.add(OP_car)
        case "cdr"    : self.compileArgs(p, vv, 1); p.add(OP_cdr)
        case "cons"   : self.compileArgs(p, vv, 2); p.add(OP_cons)
        case "eval"   : self.compileArgs(p, vv, 1); p.add(OP_eval)
        case "quote"  : self.compileQuote(p, vv)
        case "begin"  : self.compileBegin(p, vv)
        case Lambda   : fallthrough
        case "lambda" : self.compileLambda(p, vv)
        case "define" : self.compileDefine(p, vv)
        case "if"     : self.compileCondition(p, v)
        case "do"     : self.compileList(p, self.desugarDo(vv))
        case "let"    : self.compileList(p, self.desugarLet(vv, true))
        case "let*"   : self.compileList(p, self.desugarLet(vv, false))
        default       : p.i32(OP_apply, self.compileArgs(p, v, -1))
    }
}

func (self Compiler) compileArgs(p *Program, v *List, n int) uint32 {
    var nb int
    var ok bool
    var vv *List

    /* scan every element */
    for s := v; s != nil; s, nb = vv, nb + 1 {
        if vv, ok = AsList(s.Cdr); !ok {
            panic("compile: improper list is not applicable: " + v.String())
        } else {
            self.compileValue(p, s.Car)
        }
    }

    /* apply the tuple */
    if n < 0 || n == nb {
        return uint32(nb)
    } else {
        panic(fmt.Sprintf("compile: expect %d arguments, got %d.", n, nb))
    }
}

func (self Compiler) compileValue(p *Program, v Value) {
    if v.IsIdentity() {
        p.val(OP_ldconst, v)
    } else if at, ok := v.(Atom); ok {
        p.val(OP_ldvar, at)
    } else if sl, ok := AsList(v); ok {
        self.compileList(p, sl)
    } else {
        panic("fatal: compile: invalid value type")
    }
}

func (self Compiler) compileBegin(p *Program, v *List) {
    for v != nil {
        ok := false
        self.compileValue(p, v.Car)

        /* drop the value if not the last one */
        if v.Cdr != nil {
            p.add(OP_drop)
        }

        /* check for proper list */
        if v, ok = AsList(v.Cdr); !ok {
            panic("compile: begin block must be a proper list: " + v.String())
        }
    }
}

func (self Compiler) compileQuote(p *Program, v *List) {
    if v.Cdr == nil {
        p.val(OP_ldconst, v.Car)
    } else {
        panic("compile: `quote` takes exact 1 argument: " + v.String())
    }
}

func (self Compiler) compileLambda(p *Program, v *List) {

}

func (self Compiler) compileDefine(p *Program, v *List) {

}

func (self Compiler) compileCondition(p *Program, v *List) {

}

func (self Compiler) desugarDo(v *List) *List {
    return v
}

func (self Compiler) desugarLet(v *List, isParallel bool) *List {
    var decl *List
    var body *List
    var defs []Value
    var init []Value

    /* list header */
    p := v
    n := 0
    ok := false

    /* deconstruct the list */
    if decl, ok = AsList(p.Car); !ok { panic("compile: malformed let construct: " + v.String()) }
    if p   , ok = AsList(p.Cdr); !ok { panic("compile: malformed let construct: " + v.String()) }
    if body, ok = AsList(p.Car); !ok { panic("compile: malformed let construct: " + v.String()) }
    if p.Cdr != nil                  { panic("compile: malformed let construct: " + v.String()) }

    /* parse the declarations */
    for p = decl; p != nil; n++ {
        var s Atom
        var q *List

        /* get the pair, and move to next item */
        if q, ok = AsList(p.Car); !ok { panic("compile: malformed let construct: " + v.String()) }
        if p, ok = AsList(p.Cdr); !ok { panic("compile: malformed let construct: " + v.String()) }
        if s, ok = q.Car.(Atom) ; !ok { panic("compile: malformed let construct: " + v.String()) }
        if q, ok = AsList(q.Cdr); !ok { panic("compile: malformed let construct: " + v.String()) }
        if q.Cdr != nil               { panic("compile: malformed let construct: " + v.String()) }

        /* add to initializer list */
        defs = append(defs, s)
        init = append(init, q.Car)
    }

    /* let star, evaluate one value at a time */
    if !isParallel {
        n = 1
    }

    /* rebuild as lambda */
    self.desugarLetRebuild(defs, init, &body, n)
    return body
}

func (self Compiler) desugarLetRebuild(defs []Value, init []Value, body **List, p int) {
    for i := len(defs) - p; i >= 0; i-- {
        args := MakeList(defs[i:i + p]...)
        expr := MakeList(Atom(Lambda), args, *body)
        *body = MakeList(append([]Value{expr}, init[i:i + p]...)...)
    }
}