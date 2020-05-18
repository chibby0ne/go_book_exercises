// Write the corresponding Pack function. Given a struct value, Pack should
// return a URL incorporating the parameter values from the struct.

package exercise12_11

import (
    "fmt"
    "net/http"
    "reflect"
    "strings"
    "strconv"
    "net/url"
)

func Pack(ptr interface{}) (*url.URL, error) {
    v := reflect.ValueOf(ptr).Elem() // the struct variable
    if v.Type().Kind() != reflect.Struct {
        return nil, fmt.Errorf("pack: expected a struct as parameter")
    }
    var u url.Values = url.Values{}
    for i := 0; i < v.NumField(); i++ {
        fieldInfo := v.Type().Field(i)
        tag := fieldInfo.Tag
        name := tag.Get("http")
        if name == "" {
            name = strings.ToLower(fieldInfo.Name)
        }
        switch v.Field(i).Kind() {
        case reflect.Array, reflect.Slice:
            for j := 0; j < v.Field(i).Len(); j++ {
                u.Add(name, fmt.Sprintf("%v", v.Field(i).Index(j)))
            }
        default:
            u.Add(name, fmt.Sprintf("%v", v.Field(i)))
        }
    }
    uPtr, err := url.Parse(u.Encode())
    if err != nil {
        return nil, err
    }
    return uPtr, nil

}

func Unpack(req *http.Request, ptr interface{}) error {
    if err := req.ParseForm(); err != nil {
        return err
    }

    // Build map of fields keyed by effective name
    fields := make(map[string]reflect.Value)
    v := reflect.ValueOf(ptr).Elem() // the struct variable
    for i := 0; i < v.NumField(); i++ {
        fieldInfo := v.Type().Field(i) // a reflect.StructField
        tag := fieldInfo.Tag // a reflect.StructTag
        name := tag.Get("http")
        if name == "" {
            name = strings.ToLower(fieldInfo.Name)
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
