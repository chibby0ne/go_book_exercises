// Extend the field tag notation to express parameter validity requirements. For
// example, a string might need to be a valid email address or credit-card number, and an integer might need to ba a valid US ZIP code. Modify Unpack to check these requirements.

package exercise12_12

import (
    "fmt"
    "net/http"
    "reflect"
    "strings"
    "strconv"
)

type Check func(v interface{}) error

func Unpack(req *http.Request, ptr interface{}, checks map[string]Check) error {
    if err := req.ParseForm(); err != nil {
        return err
    }

    // Build map of fields keyed by effective name
    fields := make(map[string]reflect.Value)
    checksEffectiveName := make(map[string]Check)
    v := reflect.ValueOf(ptr).Elem() // the struct variable
    for i := 0; i < v.NumField(); i++ {
        fieldInfo := v.Type().Field(i) // a reflect.StructField
        // fmt.Printf("%+v\n", fieldInfo)
        tag := fieldInfo.Tag // a reflect.StructTag
        name := tag.Get("http")
        if name == "" {
            name = strings.ToLower(fieldInfo.Name)
        }
        checkKey := tag.Get("check")
        if checkKey != "" {
            if check, ok := checks[checkKey]; ok {
                // fmt.Printf("v.Field(i): %v\n", v.Field(i))
                // fmt.Printf("v.Field(i).Interface(): %v\n", v.Field(i).Interface())
                checksEffectiveName[name] = check
                // if err := check(v.Field(i).Interface()); err != nil {
                //     return fmt.Errorf("Unpack: check for field not passed: %v", err)
                // }
            }
        }
        fields[name] = v.Field(i)
    }

    // Update struct field for each parameter in the request.
    for name, values := range req.Form {
        f := fields[name]
        if !f.IsValid() {
            continue // ignore unrecognized HTTP parameters
        }
        for _, value := range values {
            check := checksEffectiveName[name]
            if err := check(value); err != nil {
                return fmt.Errorf("Unpack: check for field not passed: %v", err)
            }
            if f.Kind() == reflect.Slice {
                elem := reflect.New(f.Type().Elem()).Elem()
                if err := populate(elem, value); err != nil {
                    return fmt.Errorf("%s: %v", name, err)
                }
                f.Set(reflect.Append(f, elem))
            } else {
                if err := populate(f, value); err != nil {
                    return fmt.Errorf("%s: %v", name, err)
                }
            }
        }
    }
    return nil
}

func populate(v reflect.Value, value string) error {
    switch v.Kind() {
    case reflect.String:
        v.SetString(value)
    case reflect.Int:
        i, err := strconv.ParseInt(value, 10, 64)
        if err != nil {
            return err
        }
        v.SetInt(i)
    case reflect.Bool:
        b, err := strconv.ParseBool(value)
        if err != nil {
            return err
        }
        v.SetBool(b)
    default:
        return fmt.Errorf("unsupported kind %s", v.Type())
    }
    return nil
}
