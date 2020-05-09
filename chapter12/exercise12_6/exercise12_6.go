// Adapt encode so that as an optimization, it does not encode a field whose
// value is the zero value of its type
package exercise12_6

import (
	"bytes"
	"fmt"
	"reflect"
)

func notCompositeType(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.String, reflect.Invalid:
		return true
	default:
		return false
	}
}

func encode(buf *bytes.Buffer, v reflect.Value, prefix, indent, currentIndent string, parentIsComposite bool) error {
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
		if notCompositeType(v.Elem().Kind()) {
			currentIndent = ""
		}
		return encode(buf, v.Elem(), prefix, indent, currentIndent, parentIsComposite)
	case reflect.Array, reflect.Slice:
		var shift string
		if !parentIsComposite {
			shift = currentIndent
		}
		fmt.Fprintf(buf, "%s[", shift)
		if v.Len() == 0 {
			fmt.Fprintf(buf, "]")
			return nil
		}
		newCurrentIndent := currentIndent + indent
		for i := 0; i < v.Len(); i++ {
			if isZeroValueOfType(v.Index(i)) {
				continue
			}
			fmt.Fprintf(buf, "\n")
			if err := encode(buf, v.Index(i), prefix, indent, newCurrentIndent, false); err != nil {
				return err
			}
			if i+1 < v.Len() {
				fmt.Fprintf(buf, ",")
			}
		}
		fmt.Fprintf(buf, "\n%s]", currentIndent)
	case reflect.Struct:
		fmt.Fprintf(buf, "{")
		if v.NumField() == 0 {
			fmt.Fprintf(buf, "}")
			return nil
		}
		newCurrentIndent := currentIndent + indent
		for i := 0; i < v.NumField(); i++ {
			if isZeroValueOfType(v.Field(i)) {
				continue
			}
			fmt.Fprintf(buf, "\n")
			// Struct's field name
			fmt.Fprintf(buf, "%s%s%q: ", prefix, newCurrentIndent, v.Type().Field(i).Name)
			// Struct's field value.
			if notCompositeType(v.Field(i).Kind()) {
				newCurrentIndent = ""
			}
			if err := encode(buf, v.Field(i), prefix, indent, newCurrentIndent, true); err != nil {
				return err
			}
			if i+1 < v.NumField() {
				fmt.Fprintf(buf, ",")
			}
			newCurrentIndent = currentIndent + indent
		}
		fmt.Fprintf(buf, "\n%s}", currentIndent)

	case reflect.Map:
		var shift string
		if !parentIsComposite {
			shift = currentIndent
		}
		fmt.Fprintf(buf, "%s{", shift)
		if len(v.MapKeys()) == 0 {
			fmt.Fprintf(buf, "}")
			return nil
		}
		newCurrentIndent := currentIndent + indent
		for i, key := range v.MapKeys() {
			if isZeroValueOfType(v.MapIndex(key)) {
				continue
			}
			fmt.Fprintf(buf, "\n")
			// Write the key
			if err := encode(buf, key, prefix, indent, newCurrentIndent, parentIsComposite); err != nil {
				return err
			}
			fmt.Fprintf(buf, ": ")
			if notCompositeType(v.MapIndex(key).Kind()) {
				newCurrentIndent = ""
			}
			// Write the value
			if err := encode(buf, v.MapIndex(key), prefix, indent, newCurrentIndent, true); err != nil {
				return err
			}
			if i+1 < len(v.MapKeys()) {
				fmt.Fprintf(buf, ",")
			}
			// restore newCurrentIndent in case the Value asociated with the key was not CompositeType
			newCurrentIndent = currentIndent + indent
		}
		fmt.Fprintf(buf, "\n%s}", currentIndent)

	case reflect.Bool:
		fmt.Fprintf(buf, "%v", v.Bool())
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%f", v.Float())
	default: // chan, func, complex, interface
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func isZeroValueOfType(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if 0 == v.Int() {
			return true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if 0 == v.Uint() {
			return true
		}
	case reflect.String:
		if "" == v.String() {
			return true
		}
	case reflect.Array, reflect.Slice:
		if 0 == v.Len() {
			return true
		}
	case reflect.Struct:
		if 0 == v.NumField() {
			return true
		}
	case reflect.Map:
		if 0 == len(v.MapKeys()) {
			return true
		}
	case reflect.Bool:
		if false == v.Bool() {
			return true
		}
	case reflect.Float32, reflect.Float64:
		if 0.0 == v.Float() {
			return true
		}
	default:
		return false
	}
	return false
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v), "", "", "", false); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v), prefix, indent, "", false); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
