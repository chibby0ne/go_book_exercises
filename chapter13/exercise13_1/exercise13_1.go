// Define a deep comparison function taht considers numbers of any type equal
// if they differ by less than one part in a billion

package exercise13_1

import (
    "reflect"
    "unsafe"
)

type comparison struct {
    x, y unsafe.Pointer
    t reflect.Type
}

func equalNumbers(x, y float64) bool {
    if x == y {
        return true
    }
    var diff float64
    if x > y {
        diff = x - y
    } else {
        diff = y - x
    }
    diffPerBillion := diff * 1e9
    if diffPerBillion > x || diffPerBillion > y {
        return false
    }
    return true
}

// Equal reports whether x and y are deeply equal
func Equal(x, y interface{}) bool {
    seen := make(map[comparison]bool)
    return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}

func equal(x, y reflect.Value, seen map[comparison]bool) bool {
    if !x.IsValid() || !y.IsValid() {
        return x.IsValid() == y.IsValid()
    }
    if x.Type() != y.Type() {
        return false
    }
    if x.CanAddr() && y.CanAddr() {
        xptr := unsafe.Pointer(x.UnsafeAddr())
        yptr := unsafe.Pointer(y.UnsafeAddr())
        if xptr == yptr {
            return true // identical references
        }
        c := comparison{xptr, yptr, x.Type()}
        if seen[c] {
            return true
        }
        seen[c] = true
    }
    switch x.Kind() {

    case reflect.Bool:
        return x.Bool() == y.Bool()

    case reflect.String:
        return x.String() == y.String()

    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return equalNumbers(float64(x.Int()), float64(y.Int()))

    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return equalNumbers(float64(x.Uint()), float64(y.Uint()))

    case reflect.Float32, reflect.Float64:
        return equalNumbers(x.Float(), y.Float())

    case reflect.Chan, reflect.UnsafePointer, reflect.Func:
        return x.Pointer() == y.Pointer()

    case reflect.Ptr, reflect.Interface:
        return equal(x.Elem(), y.Elem(), seen)

    case reflect.Array, reflect.Slice:
        if x.Len() != y.Len() {
            return false
        }
        for i := 0; i < x.Len(); i++ {
            if !equal(x.Index(i), y.Index(i), seen) {
                return false
            }
        }
        return true

    // struct and map cases omitted for brevity
    case reflect.Struct:
        for i := 0; i < x.NumField(); i++ {
            if !equal(x.Field(i), y.Field(i), seen) {
                return false
            }
        }
        return true
    case reflect.Map:
        if x.Len() != y.Len() {
            return false
        }
        for _, key := range x.MapKeys() {
            if !equal(x.MapIndex(key), y.MapIndex(key), seen) {
                return false
            }
        }
        return true
    }
    panic("unreachable")
}


