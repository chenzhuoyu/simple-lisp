package main

const (
    MaxTracing = 1000000
)

func isTailCall(p Program, i int) bool {
    var ctr int
    var ins Instr

    /* trace the program to find an adjcent return, skipping all the "goto" instructions */
    for ctr = MaxTracing; ctr != 0 && i < len(p); ctr-- {
        switch ins = p[i]; ins.Op() {
            case OP_goto   : i = int(ins.Iv())
            case OP_return : return true
            default        : return false
        }
    }

    /* found infinite loops */
    return false
}

func OptimizeTailCall(p Program) {
    for i := 0; i < len(p); i++ {
        if p[i].Op() == OP_apply {
            if isTailCall(p, i + 1) {
                p[i].u0 = uint32(OP_tailcall)
            }
        }
    }
}
