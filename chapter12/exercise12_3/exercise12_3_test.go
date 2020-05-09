package exercise12_3_test

import (
	"github.com/chibby0ne/go_book_exercises/chapter12/exercise12_3"
	"testing"
)

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func TestMarshal(t *testing.T) {
	t.SkipNow()
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		// Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove":           "Peter Sellers",
			"Grp. Capt Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":      "Peter Sellers",
			"Gen. Buck Turgidson":       "George C. Scott",
			"Brig. Gen. Jack D. Ripper": "Sterlin Hayden",
			`Maj. T.J. "King" Kong`:     "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomim)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	output, err := exercise12_3.Marshal(strangelove)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(output))
	}

}

type NilInterface interface{}

type Something struct {
	NilInterface
}

func TestOtherTypes(t *testing.T) {

	something := Something{3}

	tests := []struct {
		input interface{}
		want  string
	}{

		{float32(32), `32.000000`},
		{float64(64), `64.000000`},
		{complex(float32(64), float32(-64)), `#C(64 -64)`},
		{complex(float64(128), float64(-128)), `#C(128 -128)`},
		{something, `((NilInterface ("exercise12_3_test.NilInterface" 3)))`},
	}

	var err error
	var output []byte

	for _, test := range tests {
		output, err = exercise12_3.Marshal(test.input)
		if err != nil {
			t.Error(err)
		} else if string(output) != test.want {
			t.Errorf("Marshal(%v) = %v, want %v", test.input, string(output), test.want)
		} else {
			t.Log(string(output))
		}
	}
}
