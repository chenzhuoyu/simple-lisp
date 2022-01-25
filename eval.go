package main

import (
    `fmt`
)

type Scope struct {
    prev *Scope
    defs map[string]Value
}

func CreateGlobalScope() (ret *Scope) {
    ret = new(Scope)
    ret.defs = make(map[string]Value, 16)
    ret.initAsGlobal()
    return
}

func (self *Scope) Get(key string) (v Value, ok bool) {
    for p := self; !ok && p != nil; p = p.prev { v, ok = p.defs[key] }
    return
}

func (self *Scope) Set(key string, val Value) {
    self.defs[key] = val
}

func (self *Scope) Merge(proc *Proc, vals []Value) {
    argv := len(vals)
    argc := len(proc.Args)

    /* check for args */
    if argv != argc {
        panic(fmt.Sprintf("eval: proc %s takes %d arguments, got %d", proc.Name, argc, argv))
    }

    /* fill each args */
    for i, v := range proc.Args {
        self.Set(v, vals[i])
    }
}

func (self *Scope) Derive(proc *Proc, vals []Value) (ret *Scope) {
    ret = new(Scope)
    ret.prev = self
    ret.defs = make(map[string]Value, len(proc.Args))
    ret.Merge(proc, vals)
    return
}

func (self *Scope) resolve(name string) _ValueRef {
    var ok bool
    var sc *Scope

    /* find the value */
    for sc = self; sc != nil; sc = sc.prev {
        if _, ok = sc.defs[name]; ok {
            break
        }
    }

    /* build the value reference */
    return _ValueRef {
        refs: sc,
        name: name,
    }
}

func (self *Scope) initAsGlobal() {
    for k, v := range intrinsicsTab {
        self.Set(k, v)
    }
}

type _ValueRef struct {
    refs *Scope
    name string
}

func (self _ValueRef) update(v Value) {
    if self.refs == nil {
        panic("eval: undefined reference: " + self.name)
    } else {
        self.refs.defs[self.name] = v
    }
}

func sttop(st []Value) Value {
    if nb := len(st); nb == 0 {
        panic("fatal: stack underflow")
    } else {
        return st[nb - 1]
    }
}

func stpop(st *[]Value) (v Value) {
    if i := len(*st) - 1; i < 0 {
        panic("fatal: stack underflow")
    } else {
        v, *st = (*st)[i], (*st)[:i]
        return
    }
}

func stsub(st []Value, vv Value) {
    st[len(st) - 1] = vv
}

func strem(st *[]Value, nb int) (v []Value) {
    if nb == 0 {
        panic("fatal: taking nothing")
    } else if i := len(*st) - nb; i < 0 {
        panic("fatal: stack underflow")
    } else {
        v, *st = (*st)[i:], (*st)[:i]
        return
    }
}

func istrue(v Value) bool {
    switch vv := v.(type) {
        case nil     : return false
        case Int     : return vv.BitLen() != 0
        case Frac    : return vv.Num().BitLen() != 0
        case Bool    : return bool(vv)
        case Char    : return vv != 0
        case String  : return vv != ""
        case Double  : return vv != 0.0
        case Complex : return vv != 0i
        default      : return true
    }
}

func Evaluate(s *Scope, p Program) Value {
    pc := 0
    st := make([]Value, 0, 16)

    /* execute every instruction */
    for pc < len(p) {
        iv := p[pc]
        op := iv.Op()

        /* main switch on opcode */
        switch pc++; op {
            default: {
                panic("eval: invalid instruction: " + iv.String())
            }

            /* load constant into stack */
            case OP_ldconst: {
                st = append(st, iv.Rv())
            }

            /* load proc into stack */
            case OP_ldproc: {
                st = append(st, iv.Fn().LoadWithScope(s))
            }

            /* load variable into stack */
            case OP_ldvar: {
                if vv, ok := s.Get(iv.Sv()); ok {
                    st = append(st, vv)
                } else {
                    panic("eval: undefined reference: " + iv.Sv())
                }
            }

            /* define a new variable */
            case OP_define: {
                s.Set(iv.Sv(), sttop(st))
            }

            /* set new value to an existing variable */
            case OP_set: {
                s.resolve(iv.Sv()).update(sttop(st))
            }

            /* get the first half of a pair */
            case OP_car: {
                if r, ok := sttop(st).(*List); ok {
                    stsub(st, r.Car)
                } else {
                    panic("eval: invalid argument type for car: " + AsString(sttop(st)))
                }
            }

            /* get the second half of a pair */
            case OP_cdr: {
                if r, ok := sttop(st).(*List); ok {
                    stsub(st, r.Cdr)
                } else {
                    panic("eval: invalid argument type for cdr: " + AsString(sttop(st)))
                }
            }

            /* construct a new pair from stack */
            case OP_cons: {
                cdr := stpop(&st)
                car := sttop(st)
                stsub(st, MakePair(car, cdr))
            }

            /* drop the stack top */
            case OP_drop: {
                stpop(&st)
            }

            /* unconditional jump */
            case OP_goto: {
                if pc = int(iv.Iv()); pc < 0 || pc >= len(p) {
                    panic("fatal: branch out of scope: " + iv.String())
                }
            }

            /* branch if the stack top is #f */
            case OP_if_false: {
                if !istrue(stpop(&st)) {
                    if pc = int(iv.Iv()); pc < 0 || pc >= len(p) {
                        panic("fatal: branch out of scope: " + iv.String())
                    }
                }
            }

            /* assert the stack top is true, otherwise branch */
            case OP_assert_true: {
                if istrue(sttop(st)) {
                    stpop(&st)
                } else if pc = int(iv.Iv()); pc < 0 || pc >= len(p) {
                    panic("fatal: branch out of scope: " + iv.String())
                }
            }

            /* assert the stack top is false, otherwise branch */
            case OP_assert_false: {
                if !istrue(sttop(st)) {
                    stpop(&st)
                } else if pc = int(iv.Iv()); pc < 0 || pc >= len(p) {
                    panic("fatal: branch out of scope: " + iv.String())
                }
            }

            /* apply subroutine, maybe tail-call */
            case OP_apply, OP_tailcall: {
                nb := int(iv.Iv())
                vv := strem(&st, nb)

                /* only loaded procs can be tail-called */
                if op == OP_tailcall {
                    if fn, ok := vv[0].(LoadedProc); ok {
                        if len(st) != 0 {
                            panic("fatal: unbalanced stack when tail-call")
                        } else {
                            p, pc = fn.Proc.Code, 0
                            s.Merge(fn.Proc, vv[1:])
                            break
                        }
                    }
                }

                /* check for callables */
                if fn, ok := vv[0].(Callable); !ok {
                    panic("eval: object is not appliable: " + AsString(vv[0]))
                } else {
                    st = append(st, fn.Call(vv[1:]))
                }
            }

            /* return from subroutine */
            case OP_return: {
                if len(st) != 1 {
                    panic("fatal: unbalanced stack")
                } else {
                    return st[0]
                }
            }
        }
    }

    /* should not reach here */
    dis := p.String()
    panic("fatal: program is not returned properly: \n" + dis)
}