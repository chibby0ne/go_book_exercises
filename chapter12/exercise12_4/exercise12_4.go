// Modify encode to pretty-print the S-expression in the style shown above
package exercise12_4

import (
	"bytes"
	"fmt"
	"reflect"
)

func needsExtraIndentation(k reflect.Kind) bool {
	if k == reflect.Map || k == reflect.Array || k == reflect.Slice || k == reflect.Struct {
		return true
	}
	return false
}

func encode(buf *bytes.Buffer, v reflect.Value, preindentation, postindentation int) error {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Fprintf(buf, "%v", nil)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem(), preindentation, postindentation)
	case reflect.Array, reflect.Slice:
		fmt.Fprintf(buf, "(")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				// we add one more to the pre and post indentation for the previously put (
				fmt.Fprintf(buf, "\n%*v", preindentation+1+postindentation, "")
			}
			if needsExtraIndentation(v.Index(i).Kind()) {
				postindentation += len(v.Type().Name())
			}
			if err := encode(buf, v.Index(i), preindentation, postindentation); err != nil {
				return err
			}
			if needsExtraIndentation(v.Index(i).Kind()) {
				postindentation -= len(v.Type().Name())
			}
		}
		buf.WriteByte(')')
	case reflect.Struct:
		preindentation += 1
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				fmt.Fprintf(buf, "\n%*v", preindentation, "")
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if needsExtraIndentation(v.Field(i).Kind()) {
				postindentation += len(v.Type().Field(i).Name) + 2
			}
			if err := encode(buf, v.Field(i), preindentation, postindentation); err != nil {
				return err
			}
			if needsExtraIndentation(v.Field(i).Kind()) {
				postindentation -= (len(v.Type().Field(i).Name) + 2)
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
		preindentation -= 1
	case reflect.Map:
		fmt.Fprintf(buf, "(")
		// Need to increase preindentation due to the extra ( that surrounds all the map
		// but only for the i > 0, elements
		for i, key := range v.MapKeys() {
			if i > 0 {
				// we add one more to the pre and post indentation for the previously put (
				fmt.Fprintf(buf, "\n%*v(", preindentation+1+postindentation, "")
			} else {
				fmt.Fprintf(buf, "(")
			}
			// Write the key
			if err := encode(buf, key, preindentation, postindentation); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if needsExtraIndentation(v.MapIndex(key).Kind()) {
				preindentation += len(key.Type().Name())
			}
			// Write the value
			if err := encode(buf, v.MapIndex(key), preindentation, postindentation); err != nil {
				return err
			}
			if needsExtraIndentation(v.MapIndex(key).Kind()) {
				postindentation -= len(key.Type().Name())
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Bool:
		if v.Bool() == true {
			fmt.Fprintf(buf, "%v", "t")
		} else {
			fmt.Fprintf(buf, "%v", "nil")
		}
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%f", v.Float())
	case reflect.Complex64, reflect.Complex128:
		comp := v.Complex()
		fmt.Fprintf(buf, "#C(%g %g)", real(comp), imag(comp))
	case reflect.Interface:
		buf.WriteByte('(')
		fmt.Fprintf(buf, "%q ", v.Type())
		if err := encode(buf, v.Elem(), preindentation, postindentation); err != nil {
			return err
		}
		buf.WriteByte(')')
	default: // chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v), 0, 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
