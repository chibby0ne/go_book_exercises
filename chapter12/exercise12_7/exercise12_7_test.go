package sexpr

import (
	"bytes"
	"testing"
)

type Person struct {
	Name, LastName string
	Age            int
	Address        Address
}

type Address struct {
	StreetName string
	Number     int
	PostCode   int
}

func TestDecode(t *testing.T) {
	b := []byte(`
    ((Name "Antonio")
     (LastName "Gutierrez")
     (Age 20)
     (Address ((StreetName "Andante")
              (Number 10)
              (PostCode 88818))))`)
	r := bytes.NewReader(b)
	decoder := NewDecoder(r)
	var person Person
	err := decoder.Decode(&person)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("%+v\n", person)
}
