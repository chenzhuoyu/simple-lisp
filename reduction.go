package main

func reduceSequential(vals []Value, iter func(Value, Value) Value) (ret Value) {
    ret = vals[0]
    for _, v := range vals[1:] { ret = iter(ret, v) }
    return
}

func reduceConvolution(vals []Value, test func(Value, Value) bool) bool {
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
