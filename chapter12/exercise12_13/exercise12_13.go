// Modify the S-expression encoder (§12.4) and decoder (§12.6) so that they honr the sexpr:"..." field tag in a similar manner to encoding/json.
package exercise12_13

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

func fieldByEffectiveName(name string, st reflect.Value) reflect.Value {
    fmt.Printf("name: %v, st: %v\n", name, st)
    emptyValue := reflect.Value{}
    f := st.FieldByName(name)
    if f != emptyValue {
        return f
    }
    typ := st.Type()
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        if field.Tag.Get("sexpr") == name {
            return st.FieldByName(field.Name)
        }
    }
    return emptyValue
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
            actualValue := fieldByEffectiveName(name, v)
            read(lex, actualValue)
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

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())
    case reflect.Float32, reflect.Float64:
        fmt.Fprintf(buf, "%g", v.Float())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem())
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Struct:
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
            name := v.Type().Field(i).Tag.Get("sexpr")
            if name == "" {
                name = v.Type().Field(i).Name
            }
			fmt.Fprintf(buf, "(%s ", name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Map:
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	default: // float, compelx, bool, chan, func, interface
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
