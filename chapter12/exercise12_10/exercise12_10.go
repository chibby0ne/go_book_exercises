// Extend sexpr.Unmarshal to handle the booleans, floating-point numbers, and
// interfaces encoded by your solution to Exercise 12.3. (Hint: to decode
// interfaces, you will need a mapping from the name of each supported type to
// its reflect.Type.)

package exercise12_10

import (
    "bytes"
    "fmt"
    "io"
    "reflect"
    "strconv"
    "text/scanner"
)

type Symbol string
type String string
type Int int
type StartList struct{}
type EndList struct{}

type Token interface{}

type lexer struct {
    scan scanner.Scanner
    token rune
}

func (lex *lexer) next() { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }
func (lex *lexer) consume(want rune) {
    if lex.token != want {
        panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
    }
    lex.next()
}

type Decoder struct {
    lex *lexer
    listDepth int
}

func NewDecoder(r io.Reader) *Decoder {
    s := scanner.Scanner{}
    ss := s.Init(r)
    return &Decoder{lex: &lexer{scan: *ss}}
}

func (d *Decoder) Token() (Token, error) {
    d.lex.next()
    if d.listDepth == 0 && d.lex.token != '(' && d.lex.token != scanner.EOF {
        return nil, fmt.Errorf("sexpr: expected: '(' but got: %s", d.lex.text())
    }
    switch d.lex.token {
    case scanner.EOF:
        return nil, io.EOF
    case scanner.Ident:
        return Symbol(d.lex.text()), nil
    case scanner.String:
        text := d.lex.text()
        s, err := strconv.Unquote(text)
        if err != nil {
            return nil, fmt.Errorf("sexpr: identified a string but it is not quoted: %s", text)
        }
        return String(s), nil
    case scanner.Int:
        text := d.lex.text()
        i, err := strconv.Atoi(text)
        if err != nil {
            return nil, fmt.Errorf("sexpr: identified an int but coult not be parsed: %s", text)
        }
        return Int(i), nil
    case '(':
        d.listDepth++
        return StartList{}, nil
    case ')':
        d.listDepth--
        return EndList{}, nil
    default:
        pos := d.lex.scan.Pos()
        return nil, fmt.Errorf("sexpr: unexpected token: %s at line: %d, column: %d",  d.lex.text(), pos.Line, pos.Line)
    }
}

func Unmarshal(data []byte, out interface{}) (err error) {
    decoder := NewDecoder(bytes.NewReader(data))
    decoder.lex.next()
    defer func() {
        if x := recover(); x != nil {
            err = fmt.Errorf("error at %s: %v", decoder.lex.scan.Position, x)
        }
    }()
    read(decoder.lex, reflect.ValueOf(out).Elem())
    return nil
}

func read(lex *lexer, v reflect.Value){
    switch lex.token {
    case scanner.Ident:
        if lex.text() == "nil" {
            v.Set(reflect.Zero(v.Type()))
            lex.next()
            return
        } else if lex.text() == "t" {
            v.SetBool(true)
            lex.next()
            return
        }
    case scanner.String:
        s, _ := strconv.Unquote(lex.text())
        v.SetString(s)
        lex.next()
        return
    case scanner.Int:
        i, _ := strconv.Atoi(lex.text())
        v.SetInt(int64(i))
        lex.next()
        return
    case scanner.Float:
        f, _ := strconv.ParseFloat(lex.text(), 64)
        v.SetFloat(f)
        lex.next()
        return
    case '(':
        lex.next()
        readList(lex, v)
        lex.next()
        return
    }
    panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) {
    switch v.Kind() {
    case reflect.Array: // (item ...)
        for i := 0; !endList(lex); i++ {
            read(lex, v.Index(i))
        }
    case reflect.Slice: // (item ...)
        for !endList(lex) {
            item := reflect.New(v.Type().Elem()).Elem()
            read(lex, item)
            v.Set(reflect.Append(v, item))
        }
    case reflect.Struct: // ((name value) ...)
        for !endList(lex) {
            lex.consume('(')
            if lex.token != scanner.Ident {
                panic(fmt.Sprintf("got token %q, want field name", lex.text()))
            }
            name := lex.text()
            lex.next()
            read(lex, v.FieldByName(name))
            lex.consume(')')
        }
    case reflect.Map: // ((key value) ...)
        v.Set(reflect.MakeMap(v.Type()))
        for !endList(lex) {
            lex.consume('(')
            key := reflect.New(v.Type().Key()).Elem()
            read(lex, key)
            value := reflect.New(v.Type().Elem()).Elem()
            read(lex, value)
            v.SetMapIndex(key, value)
            lex.consume(')')
        }
    case reflect.Interface:
        name := lex.text()
        typ, ok := KnownInterfaces[name]
        if !ok {
            panic(fmt.Errorf("Interface: %q not found", name))
        }
        newVal := reflect.New(typ)
        lex.next()
        read(lex, newVal.Elem())
        v.Set(newVal.Elem())
    default:
        panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
    }
}

func endList(lex *lexer) bool {
    switch lex.token {
    case scanner.EOF:
        panic("end of file")
    case ')':
        return true
    }
    return false
}

var KnownInterfaces map[string]reflect.Type

func init() {
    KnownInterfaces = make(map[string]reflect.Type)

}
