package exercise12_10

import (
    "testing"
    "reflect"
    "strings"
)

type Person struct {
    Name string
    Age int
    Height float64
}

func TestUnmarshalStruct(t *testing.T) {
    tests := []struct{
        input string
        want interface{}
    }{
        {
        `((Name "string") (Age 19) (Height 1.88))`,
        Person{Name: "string", Age: 19, Height: 1.88},
        },
    }
    for _, test := range tests {
        var p Person
        err := Unmarshal([]byte(test.input), &p)
        if err != nil {
            t.Error(err)
        }
        if !reflect.DeepEqual(p, test.want) {
            t.Errorf("not equal: got %v, want %v", p, test.want)
        }


    }
}
// type Yeller interface{
//     Yell(msg string) string
// }

type yeller string

func (y *yeller) Yell(msg string) string {
    return strings.ToUpper(msg)
}

type Yeller interface{}

func TestUnmarshalInterface(t *testing.T) {
    tests := []struct{
        input string
        want interface{}
    }{
        {
        `(Yeller "hey")`,
        Yeller("hey"),
        },
    }
    KnownInterfaces["Yeller"] = reflect.TypeOf(string(""))
    for _, test := range tests {
        var yy Yeller
        err := Unmarshal([]byte(test.input), &yy)
        if err != nil {
            t.Error(err)
        }
        if !reflect.DeepEqual(yy, test.want) {
            t.Errorf("not equal: got %v, want %v", yy, test.want)
        }
    }
}
