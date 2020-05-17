// Create a streaming API for the S-expression encoder, following the style of json.Encoder

package sexpr

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"text/scanner"
)

type Decoder struct {
	scanner      scanner.Scanner
	currentToken rune
}

func NewDecoder(r io.Reader) *Decoder {
	s := scanner.Scanner{}
	ss := s.Init(r)
	d := &Decoder{
		scanner: *ss,
	}
	d.next()
	return d
}

func (d *Decoder) next() {
	d.currentToken = d.scanner.Scan()
}

func (d *Decoder) tokenText() string {
	return d.scanner.TokenText()
}

func (d *Decoder) discard(want rune) error {
	if d.currentToken != want {
		return fmt.Errorf("sexpr: invalid character '%c', but expected a '%c' for struct", d.currentToken, want)
	}
	d.next()
	return nil
}

func (d *Decoder) Decode(v interface{}) error {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr || vv.IsNil() {
		return fmt.Errorf("sexpr: Unmarshal(non-pointer %v)", vv.Type())
	}
	return unmarshal(d, vv.Elem())
}

func unmarshal(d *Decoder, v reflect.Value) error {
	switch d.currentToken {
	case scanner.Int:
		i, err := strconv.Atoi(d.tokenText())
		if err != nil {
			return fmt.Errorf("sexpr: Converting to int: %v", err)
		}
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(i))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(i))
		default:
			return fmt.Errorf("sexpr: expected int, but type of interface is: %v", v.Kind())
		}
		d.next()
	case scanner.Float:
		switch v.Kind() {
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(d.tokenText(), 64)
			if err != nil {
				return fmt.Errorf("sexpr: Converting to float: %v", err)
			}
			v.SetFloat(f)
		default:
			return fmt.Errorf("sexpr: expected int, but type of interface is: %v", v.Kind())
		}
		d.next()
	case scanner.Ident:
		x := d.tokenText()
		if x == "nil" {
			v.Set(reflect.Zero(v.Type()))
			d.next()
		}
	case scanner.String:
		s, err := strconv.Unquote(d.tokenText())
		if err != nil {
			return fmt.Errorf("sexpr: unquoting string: %v", err)
		}
		v.SetString(s)
		d.next()
	case '(':
		d.next()
		unmarshalList(d, v)
		d.next()
	default:
		return fmt.Errorf("Not supported type for token: %v", d.currentToken)
	}
	return nil
}

func unmarshalList(d *Decoder, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Struct:
		for end, err := endList(d.currentToken, "struct"); !end; {
			if err != nil {
				return err
			}
			if err := d.discard('('); err != nil {
				return err
			}
			if d.currentToken != scanner.Ident {
				return fmt.Errorf("sexpr: expected ident, but got: %v", d.currentToken)
			}
			name := d.tokenText()
			d.next()
			if err := unmarshal(d, v.FieldByName(name)); err != nil {
				return err
			}
			if err := d.discard(')'); err != nil {
				return err
			}
		}
	case reflect.Map:
		v.Set(reflect.MakeMap(v.Type()))
		for end, err := endList(d.currentToken, "map"); !end; {
			if err != nil {
				return err
			}
			if err := d.discard('('); err != nil {
				return err
			}
			if d.currentToken == scanner.Ident {
				return fmt.Errorf("sexpr: expected key, but got ident: %v", d.currentToken)
			}
			key := reflect.New(v.Type().Key()).Elem()
			if err := unmarshal(d, key); err != nil {
				return err
			}
			if d.currentToken == scanner.Ident {
				return fmt.Errorf("sexpr: expected value, but got ident: %v", d.currentToken)
			}
			value := reflect.New(v.Type().Elem()).Elem()
			if err := unmarshal(d, value); err != nil {
				return err
			}
			v.SetMapIndex(key, value)
			if err := d.discard(')'); err != nil {
				return err
			}
		}
	case reflect.Array:
		for i := 0; ; i++ {
			end, err := endList(d.currentToken, "array")
			if err != nil {
				return err
			}
			if end {
				return nil
			}
			if err := unmarshal(d, v.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; ; i++ {
			end, err := endList(d.currentToken, "slice")
			if err != nil {
				return err
			}
			if end {
				return nil
			}
			elem := reflect.New(v.Type().Elem()).Elem()
			if err = unmarshal(d, elem); err != nil {
				return err
			}
			v.Set(reflect.Append(v, elem))
		}

	default:
		return fmt.Errorf("Not supported type for interface: %v", v.Kind())
	}
	return nil
}

func endList(cur rune, t string) (bool, error) {
	switch cur {
	case scanner.EOF:
		return false, fmt.Errorf("sexpr: %s not correctly terminated", t)
	case ')':
		return true, nil
	}
	return false, nil
}
