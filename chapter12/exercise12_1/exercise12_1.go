// Extend Display so that it can display maps whose keys are structs or arrays
package exercise12_1

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Any(value interface{}) string {
	return formatAtom(reflect.ValueOf(value))
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		com := v.Complex()
		r := strconv.FormatFloat(real(com), 'f', -1, 64)
		i := strconv.FormatFloat(imag(com), 'f', -1, 64)
		return r + i + "i"
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default:
		return v.Type().String() + "Value"
	}
}

func formatMapKeys(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Struct:
		var s strings.Builder
		s.WriteString(fmt.Sprintf("%v {", v.Type().Name()))
		if v.NumField() != 0 {
			s.WriteString(formatMapKeys(v.Field(0)))
		}
		for i := 1; i < v.NumField(); i++ {
			s.WriteString(" " + formatMapKeys(v.Field(i)))
		}
		s.WriteRune('}')
		return s.String()
	case reflect.Array:
		var s strings.Builder
		s.WriteRune('[')
		if v.Len() != 0 {
			s.WriteString(formatMapKeys(v.Index(0)))
		}
		for i := 1; i < v.Len(); i++ {
			s.WriteString(" " + formatMapKeys(v.Index(i)))
		}
		s.WriteRune(']')
		return s.String()
	case reflect.Interface:
		if v.IsNil() {
			return fmt.Sprintf("%v %v", v.Type().Name(), fmt.Sprintf("%v", nil))
		} else {
			return fmt.Sprintf("%v %v %v", v.Type().Name(), v.Elem().Type(), v.Elem())
		}
	default:
		return formatAtom(v)
	}
}

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x))
}

func display(path string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i))
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatMapKeys(key)), v.MapIndex(key))
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem())
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}

}
