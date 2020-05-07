package display_test

import (
	"github.com/chibby0ne/go_book_exercises/chapter12/display"
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

func TestDisplayStruct(t *testing.T) {

	strangelove := Movie{
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

	display.Display("strangelove", strangelove)
}

func TestDisplayEmptyInterface(t *testing.T) {

	var i interface{} = 3

	display.Display("i", i)
	display.Display("&i", &i)
}
