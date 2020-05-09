// Adapt encode to emit JSON instead of S-expressions. Test your encoder using the standard decoder, json.Unmarshal
package exercise12_5

import (
	"bytes"
	"fmt"
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
		return encode(buf, v.Elem(), prefix, indent, currentIndent, parentIsComposite)
	case reflect.Array, reflect.Slice:
		newCurrentIndent := currentIndent + indent
		if parentIsComposite {
			fmt.Fprintf(buf, "[")
		} else {
			fmt.Fprintf(buf, "%s[", currentIndent)
		}
		for i := 0; i < v.Len(); i++ {
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
		newCurrentIndent := currentIndent + indent

		for i := 0; i < v.NumField(); i++ {
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
		if parentIsComposite {
			fmt.Fprintf(buf, "{")
		} else {
			fmt.Fprintf(buf, "%s{", currentIndent)
		}
		newCurrentIndent := currentIndent + indent
		for i, key := range v.MapKeys() {
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
