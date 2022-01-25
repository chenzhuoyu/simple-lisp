package main

func reduce(vals []Value, iter func(Value, Value) Value) (ret Value) {
    ret = vals[0]
    for _, v := range vals[1:] { ret = iter(ret, v) }
    return
}

func reduceSeq(vals []Value, defv Value, iter func(Value, Value) Value) (ret Value) {
    switch len(vals) {
        case 0  : return defv
        case 1  : return iter(defv, vals[0])
        case 2  : return iter(vals[0], vals[1])
        default : return reduce(vals, iter)
    }
}

func reduceMonotonic(vals []Value, test func(Value, Value) bool) bool {
    i := 0
    r := true
    n := len(vals)

    /* fold through the list */
    for r && i < n - 1 {
        if i++; !test(vals[i - 1], vals[i]) {
            r = false
        }
    }

    /* all done */
    return r
}
