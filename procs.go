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

func (self *Proc) LoadWithScope(scope *Scope) LoadedProc {
    return LoadedProc {
        Proc  : self,
        Scope : scope,
    }
}

type LoadedProc struct {
    *Proc
    *Scope
}

func (self LoadedProc) Call(args []Value) Value {
    return Evaluate(self.Scope.Derive(self.Proc, args), self.Code)
}
