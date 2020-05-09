package exercise12_5_test

import (
	"encoding/json"
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter12/exercise12_5"
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

func equalMovies(a, b *Movie) bool {
	if a.Title != b.Title {
		return false
	}
	if a.Subtitle != b.Subtitle {
		return false
	}
	if a.Year != b.Year {
		return false
	}
	if a.Color != b.Color {
		return false
	}
	for k, v := range a.Actor {
		v2, ok := b.Actor[k]
		if !ok {
			return false
		}
		if v != v2 {
			return false
		}
	}
	for i, v := range a.Oscars {
		if v != b.Oscars[i] {
			return false
		}
	}
	if a.Sequel != b.Sequel {
		return false
	}
	return true

}

var strangelove = Movie{
	Title:    "Dr. Strangelove",
	Subtitle: "How I learned to Stop Worrying and Love the Bomb",
	Year:     1964,
	Color:    false,
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

func TestMarshal(t *testing.T) {
	tests := []struct {
		input interface{}
	}{
		{
			strangelove,
		},
		{
			[]int{1, 2, 3},
		},
		{
			map[string]string{
				"hey": "there",
				"how": "is it going?",
			},
		},
		{float64(3.0)},
		{interface{}(3.0)},
		{},

		{
			map[string]map[string]string{
				"name": {
					"real":  "joe",
					"alias": "the machine",
				},
				"location": {
					"country":   "usa",
					"continent": "north america",
				},
			},
		},
		{
			[][]int{
				{1, 2, 3},
				{10, 20, 30},
			},
		},
	}

	for _, test := range tests {
		output, err := exercise12_5.MarshalIndent(test.input, "", "  ")
		if err != nil {
			t.Error(err)
		}
		goldenOutput, err := json.MarshalIndent(test.input, "", "  ")
		if err != nil {
			t.Error(err)
		}

		fmt.Println(string(output))
		fmt.Println(string(goldenOutput))

	}

}

func TestMarshalUnmarshal(t *testing.T) {
	output, err := exercise12_5.MarshalIndent(strangelove, "", "  ")
	if err != nil {
		t.Error(err)
	}

	var movie Movie
	if err := json.Unmarshal(output, &movie); err != nil {
		t.Error(err)
	}

	if !equalMovies(&strangelove, &movie) {
		t.Errorf("got: %v, want: %v", movie, strangelove)
	}
}

func equalComplex(a, b *complex128) bool {
	if real(*a) != real(*b) || imag(*a) != imag(*b) {
		return false
	}
	return true
}

// func TestMarshalUnmarshalComplex(t *testing.T) {
// 	compl := complex(float64(64), float64(-64))
// 	goldenOutput, err := json.MarshalIndent(compl, "", "  ")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	output, err := exercise12_5.MarshalIndent(compl, "", "  ")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	var c, d complex128
// 	if err := json.Unmarshal(output, &c); err != nil {
// 		t.Error(err)
// 	}
// 	if err := json.Unmarshal(goldenOutput, &d); err != nil {
// 		t.Error(err)
// 	}
// 	if !equalComplex(&compl, &c) {
// 		t.Errorf("got: %v, want: %v\n should have gotten: %v, should want: %v", string(output), c, string(goldenOutput), d)
// 	}
// }
