package exercise12_11

import (
    "testing"
)


type data struct{
    Labels []string `http:"l"`
    MaxResults int `http:"max"`
    Exact bool `http:"x"`
}

type Auth struct{
    Key string `http:"key"`
    AgentName string `http:"name"`
    ExpirationDate int `http:"expires"`
}

func TestPack(t *testing.T) {
    tests := []struct {
        input interface{}
        want string
    }{
        {
            data{
                Labels: []string{"golang", "programming"},
                MaxResults: 20,
                Exact: true,

            },
            "l=golang&l=programming&max=20&x=true",
        },
        {
            Auth{
                Key: "abracadabra",
                AgentName: "coulson",
                ExpirationDate: 100,
            },
            "expires=100&key=abracadabra&name=coulson",
        },
    }
    for _, test := range tests {
        switch input := test.input.(type) {
        case data:
            url, err := Pack(&input)
            if err != nil {
                t.Error(err)
            } else if url.String() != test.want {
                t.Errorf("got: %v, want: %v", url.String(), test.want)
            }
        case Auth:
            url, err := Pack(&input)
            if err != nil {
                t.Error(err)
            } else if url.String() != test.want {
                t.Errorf("got: %v, want: %v", url.String(), test.want)
            }
        default:
            panic("unexpected input type")
        }
        
    }
}


