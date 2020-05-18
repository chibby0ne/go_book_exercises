package exercise12_12

import (
    "fmt"
    "net/http"
    "testing"
    "net/url"
    "reflect"
)

func CheckValidUSZipCode(v interface{}) error {
    fmt.Println(v)
    val, ok := v.(string)
    if !ok {
        return fmt.Errorf("field value: %v is not a string", val)
    }
    // All US ZipCodes are 5 digits long
    if len(val) != 5 {
        return fmt.Errorf("field is not a valid zip code, it's length is not 5 digits: %v", val)
    }
    return nil
}

type data struct {
    ZipCode string `http:"zip" check:"validzipcode"`
}

func TestUnpack(t *testing.T) {
    tests := []struct {
        req *http.Request
        ptr data
        err error
    }{
        {
            &http.Request{
                Form: url.Values{
                    "http": []string{"zip"},
                    "check": []string{"validzipcode"},
                },
            },
            data{
                ZipCode: "01234",
            },
            nil,
        },
        {
            &http.Request{
                Form: url.Values{
                    "http": []string{"zip"},
                    "check": []string{"validzipcode"},
                },
            },
            data{
                ZipCode: "101234",
            },
            fmt.Errorf(fmt.Sprintf("Unpack: check for field not passed: %v", fmt.Errorf("field is not a valid zip code, it's length is not 5 digits: %v", "101234"))),

        },

    }
    checks := map[string]Check{
        "validzipcode": CheckValidUSZipCode,
    }

    for _, test := range tests {
        var d data
        err := Unpack(test.req, &d, checks)
        if err != nil && test.err != nil {
            if err != test.err {
                t.Errorf("got: %v, want: %v", err, test.err)
            }
        }
        if reflect.DeepEqual(d, test.ptr) {
            t.Errorf("got: %v, want: %v", d, test.ptr)
        }
    }
    
}


