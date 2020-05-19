// Write a function that reports whether its argument is a cyclic data
// structure

package exercise13_2

import (
    "reflect"
    "unsafe"
    // "fmt"
)

type comparison struct {
    x, y unsafe.Pointer
    t reflect.Type
}

func ReportCyclicStruct(x interface{}) bool {
    v := reflect.ValueOf(x)
    t := reflect.TypeOf(x)
    // fmt.Printf("from ReportCyclicStruct: v =  %+v\n", v)
    switch v.Kind() {
    case reflect.Struct:
        // fmt.Printf("v  = %+v, is a struct\n", v)
    case reflect.Ptr:
        // fmt.Printf("v  = %+v, is a pointer\n", v)
        if v.Elem().Kind() != reflect.Struct {
            return false
        }
        v = v.Elem()
        t = v.Type()
        // fmt.Printf("v is now = %+v\n", v)
        // fmt.Printf("t is now = %+v\n", t)
    default:
        return false
    }
    if reportCyclicStruct(v, t) {
        // fmt.Println("This is a cyclic data structure")
        return true
    }
    return false
}

func reportCyclicStruct(x reflect.Value, t reflect.Type) bool {
    if x.Kind() != reflect.Struct {
        return false
    }
    // fmt.Printf("t: %v\n", t)
    for i := 0; i < x.NumField(); i++ {
        // fmt.Printf("x.Field(i).Type(): %v\n", x.Field(i).Type())
        if x.Field(i).Kind() == reflect.Ptr {
            if t == x.Field(i).Type().Elem() || reportCyclicStruct(x.Field(i), t) || reportCyclicStruct(x.Field(i), x.Field(i).Type()) {
                return true
            }
        }
    }
    return false
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
        // if reportCyclicStruct(x, x.Type()) {
        //     fmt.Println("This is a cyclic data structure")
        // }
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


