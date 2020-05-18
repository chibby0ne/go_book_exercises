// Exercise 12.8: The sexpr.Unmarshal function, like json.Marshal, requires the
// complete input in a byte slice before it can begin decoding. Define a
// sexpr.Decoder type that, like json.Decoder, allows a sequence of values to
// be decoded from an io.Reader. Change sexpr.Unmarshal to use this new type.
package exercise12_8

import (
    "bytes"
    "fmt"
    "io"
    "reflect"
    "strconv"
    "text/scanner"
)

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
}

func NewDecoder(r io.Reader) *Decoder {
    s := scanner.Scanner{}
    ss := s.Init(r)
    return &Decoder{&lexer{scan: *ss}}
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
