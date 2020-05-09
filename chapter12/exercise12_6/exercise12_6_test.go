package exercise12_6_test

import (
	"encoding/json"
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter12/exercise12_6"
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
		output, err := exercise12_6.MarshalIndent(test.input, "", "  ")
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

type Person struct {
	name string
	car  Car
}

type Car struct {
	model     string
	brand     string
	sold      int
	year      int
	mpg       float64
	price     uint
	extras    []int
	signature string
	brandNew  bool
}

func TestPrettyPrintNestedStruct(t *testing.T) {
	someone := Person{
		name: "joe",
		car: Car{
			model:     "model s",
			brand:     "tesla",
			year:      2019,
			signature: "",
			extras:    []int{},
			sold:      0,
			price:     0,
			mpg:       0.000,
			brandNew:  false,
		},
	}
	output, err := exercise12_6.MarshalIndent(someone, "", "  ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(output))
}
