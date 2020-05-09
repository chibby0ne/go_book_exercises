package exercise12_4_test

import (
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter12/exercise12_4"
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
	tests := []struct {
		input interface{}
	}{
		{
			Movie{
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
			},
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
	}

	for _, test := range tests {
		output, err := exercise12_4.Marshal(test.input)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println(string(output))
		}
	}

}
