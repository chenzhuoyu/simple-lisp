package main

import (
    `fmt`
    `strconv`
    `strings`
    `unicode`
)

var _CharTab = map[string]rune {
    "space"     : ' ',
    "newline"   : '\n',
    "backspace" : '\b',
    "tab"       : '\t',
    "linefeed"  : '\n',
    "page"      : '\f',
    "return"    : '\r',
    "rubout"    : 0x7f,
}

func isAtomChar(ch rune) bool {
    return !(ch == '(' || ch == ')' || ch == '"' || unicode.IsSpace(ch))
}

type Parser struct {
    p int
    s []rune
}

func CreateParser(src string) *Parser {
    return &Parser {
        p: 0,
        s: []rune(src),
    }
}

func (self *Parser) error(msg string) error {
    var row int
    var col int

    /* count row and coloumn */
    for _, ch := range self.s[:self.p] {
        if col++; ch == '\n' {
            row++
            col = 0
        }
    }

    /* build the error */
    return fmt.Errorf(
        "syntax error at row %d, column %d: %s",
        row + 1,
        col + 1,
        msg,
    )
}

func (self *Parser) noEOF(topLevel bool) {
    if !topLevel && self.p >= len(self.s) {
        panic(self.error("unexpected EOF"))
    }
}

func (self *Parser) noSpace() {
    for self.p < len(self.s) && unicode.IsSpace(self.s[self.p]) {
        self.p++
    }
}

func (self *Parser) nextChar() (cc rune) {
    if i := self.p; i >= len(self.s) {
        return 0
    } else {
        self.p++
        return self.s[i]
    }
}

func (self *Parser) parseStr() Value {
    q := self.p
    p := self.p - 1

    /* scan until the end of string */
    for q < len(self.s) && self.s[q] != '"' {
        if q++; self.s[q - 1] == '\\' {
            q++
        }
    }

    /* check for string termination */
    if q >= len(self.s) {
        panic(self.error("string is not terminated"))
    }

    /* unquote the string */
    self.p = q + 1
    ret, err := strconv.Unquote(string(self.s[p:self.p]))

    /* check for errors */
    if err != nil {
        panic(self.error("cannot parse string literal: " + err.Error()))
    } else {
        return String(ret)
    }
}

func (self *Parser) parseCdr() Value {
    if ret := self.parseList(false); ret == nil {
        return nil
    } else {
        return ret
    }
}

func (self *Parser) parseList(topLevel bool) *List {
    for p, q := (*List)(nil), (*List)(nil);; {
        if vv, ok := self.parseValue(topLevel); !ok || vv == Atom(")") {
            return p
        } else if vv != Atom(".") {
            AppendValue(&p, &q, vv)
        } else if q == nil {
            panic(self.error("ill-formed dotted list"))
        } else if q.Cdr, ok = self.parseValue(false); !ok {
            panic(self.error("cdr expression expected"))
        } else if vv, ok = self.parseValue(false); !ok || vv != Atom(")") {
            panic(self.error("')' expected"))
        } else {
            return p
        }
    }
}

func (self *Parser) parseChar(ch string) Char {
    if _CharTab[ch] != 0 {
        return Char(_CharTab[ch])
    } else if chars := []rune(ch); len(chars) == 1 {
        return Char(chars[0])
    } else {
        panic(self.error(`invalid character name #\` + ch))
    }
}

func (self *Parser) parseValue(topLevel bool) (Value, bool) {
    self.noSpace()
    self.noEOF(topLevel)

    /* check for simple cases */
    switch self.nextChar() {
        case 0    : return nil, false
        case '\'' : break
        case ')'  : return Atom(")"), true
        case '"'  : return self.parseStr(), true
        case '('  : return self.parseCdr(), true
        default   : return self.parseSimple(), true
    }

    /* desugar "quote" keyword */
    if v, ok := self.parseValue(false); !ok {
        return nil, false
    } else {
        return MakeList(Atom("quote"), v), true
    }
}

func (self *Parser) parseSimple() Value {
    p := self.p - 1
    n := len(self.s)

    /* scan until the next space or EOF */
    for self.p < n && isAtomChar(self.s[self.p]) {
        self.p++
    }

    /* slice the value */
    val := string(self.s[p:self.p])
    low := strings.ToLower(val)

    /* check for token types */
    if low == "#t" {
        return Bool(true)
    } else if low == "#f" {
        return Bool(false)
    } else if strings.HasPrefix(low, `#\`) {
        return self.parseChar(low[2:])
    } else if iv, err := strconv.ParseInt(val, 0, 64); err == nil {
        return Int(iv)
    } else if fv, err := strconv.ParseFloat(val, 64); err == nil {
        return Float(fv)
    } else if cv, err := strconv.ParseComplex(val, 128); err == nil {
        return Complex(cv)
    } else {
        return Atom(val)
    }
}

func (self *Parser) Parse() *List {
    return &List {
        Car: Atom("begin"),
        Cdr: self.parseList(true),
    }
}
