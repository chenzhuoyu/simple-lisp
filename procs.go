package main

import (
    `fmt`
    `strings`
)

type Proc struct {
    Name string
    Code Program
    Args []string
}

func (self *Proc) Call(sc *Scope, args []Value) Value {
    argv := len(args)
    argc := len(self.Args)
    subr := sc.derive()

    /* check for args */
    if argv != argc {
        panic(fmt.Sprintf("eval: proc %s takes %d arguments, got %d", self.Name, argc, argv))
    }

    /* fill each args */
    for i, v := range self.Args {
        subr.Set(v, args[i])
    }

    /* evaluate the program */
    return Evaluate(subr, self.Code)
}

func (self *Proc) String() string {
    if len(self.Args) == 0 {
        return fmt.Sprintf("#[proc (%s)]", self.Name)
    } else {
        return fmt.Sprintf("#[proc (%s %s)]", self.Name, strings.Join(self.Args, " "))
    }
}

func (self *Proc) IsIdentity() bool {
    return true
}
