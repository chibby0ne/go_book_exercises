// Modify encode to pretty-print the S-expression in the style shown above
package exercise12_5

import (
	"bytes"
	// "strconv"
	// "encoding/json"
	"fmt"
	// "math"
	"reflect"
)

func notCompositeType(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.String:
		return true
	default:
		return false
	}
}

func encode(buf *bytes.Buffer, v reflect.Value, prefix, indent, currentIndent string) error {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Fprintf(buf, "%s%s%v", prefix, currentIndent, "null")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%s%s%d", prefix, currentIndent, v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%s%s%d", prefix, currentIndent, v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%s%s%q", prefix, currentIndent, v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem(), prefix, indent, currentIndent)
	case reflect.Array, reflect.Slice:
		fmt.Fprintf(buf, "%s[", currentIndent)
		newindent := currentIndent + indent
		for i := 0; i < v.Len(); i++ {
			fmt.Fprintf(buf, "\n")
			if err := encode(buf, v.Index(i), prefix, newindent, newindent); err != nil {
				return err
			}
			if i+1 < v.Len() {
				fmt.Fprintf(buf, ",")
			}
		}
		fmt.Fprintf(buf, "\n%s]", currentIndent)
	case reflect.Struct:
		fmt.Fprintf(buf, "{")
		currentIndent += indent
		newindent := currentIndent
		for i := 0; i < v.NumField(); i++ {
			fmt.Fprintf(buf, "\n")
			// Struct's field name
			fmt.Fprintf(buf, "%s%s%q: ", prefix, currentIndent, v.Type().Field(i).Name)
			// Struct's field value.
			// For the special case of arrays/slices that add the currentIndent before
			// the "[" for the case when they are nested
			// if v.Field(i).Kind() == reflect.Array || v.Field(i).Kind() == reflect.Slice {
			// 	newindent = indent
			// }
			if notCompositeType(v.Field(i).Kind()) {
				newindent = ""
			}
			if err := encode(buf, v.Field(i), prefix, indent, newindent); err != nil {
				return err
			}
			if i+1 < v.NumField() {
				fmt.Fprintf(buf, ",")
			}
			newindent = currentIndent
		}
		fmt.Fprintf(buf, "\n}")
	case reflect.Map:
		fmt.Fprintf(buf, "{")
		newindent := currentIndent + indent
		for i, key := range v.MapKeys() {
			fmt.Fprintf(buf, "\n")
			// Write the key
			if err := encode(buf, key, prefix, newindent, newindent); err != nil {
				return err
			}
			fmt.Fprintf(buf, ": ")
			// Write the value
			if err := encode(buf, v.MapIndex(key), prefix, newindent, ""); err != nil {
				return err
			}
			if i+1 < len(v.MapKeys()) {
				fmt.Fprintf(buf, ",")
			}
		}
		fmt.Fprintf(buf, "\n%s}", currentIndent)
	case reflect.Bool:
		fmt.Fprintf(buf, "%v", v.Bool())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%f", v.Float())
	// case reflect.Complex64, reflect.Complex128:
	// 	comp := v.Complex()
	// 	r, i := real(comp), imag(comp)
	// 	r_f := strconv.FormatFloat(r, 'f', -1, 64)
	// 	i_f := strconv.FormatFloat(math.Abs(i), 'f', -1, 64)
	// 	var sign rune = '+'
	// 	if i < 0 {
	// 		sign = '-'
	// 	}
	// 	fmt.Fprintf(buf, "(%s%c%si)", r_f, sign, i_f)
	// case reflect.Interface:
	// 	buf.WriteByte('{')
	// 	fmt.Fprintf(buf, "%q ", v.Type())
	// 	if err := encode(buf, v.Elem(), prefix, indent, currentIndent); err != nil {
	// 		return err
	// 	}
	// 	buf.WriteByte('}')
	default: // chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v), "", "", ""); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v), prefix, indent, ""); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
