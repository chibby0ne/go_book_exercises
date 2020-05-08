package exercise12_1_test

import (
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter12/exercise12_1"
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

	exercise12_1.Display("strangelove", strangelove)

}

func TestDisplayEmptyInterface(t *testing.T) {
	var i interface{} = 3

	exercise12_1.Display("i", i)
	exercise12_1.Display("&i", &i)
}

func TestMapKeys(t *testing.T) {
	m := make(map[[5]string]string)

	array := [5]string{"hola", "como", "estas", "tu", "amigo"}
	array2 := [5]string{"hey", "there", "how", "are", "you?"}
	m[array] = "spanish"
	m[array2] = "english"

	fmt.Println("Map using array as key")
	for k, v := range m {
		fmt.Printf("m[%v] = %v\n", k, v)
	}

	exercise12_1.Display("m", m)

	mm := make(map[Person]string)

	joe := NewRegularPerson("joe")
	daphne := NewRegularPerson("daphne")
	mm[joe] = "asdfasdf"
	mm[daphne] = "koller"
	mm[nil] = "hoho"

	fmt.Println("Map using interface as key")
	for k, v := range mm {
		fmt.Printf("mm[%v] = %v\n", k, v)
	}
	fmt.Printf("nil: %v\n", nil)

	exercise12_1.Display("mm", mm)

	mmm := make(map[RegularPerson]string)
	mmm[*joe] = "rogan"
	mmm[*daphne] = "koller"

	fmt.Println("Map using struct as key")
	for k, v := range mmm {
		fmt.Printf("mmm[%v] = %v\n", k, v)
	}

	exercise12_1.Display("mmm", mmm)

}

type Person interface {
	Name() string
}

type RegularPerson struct {
	name string
	last string
}

func NewRegularPerson(name string) *RegularPerson {
	if name == "" {
		name = "charlie"
	}
	return &RegularPerson{
		name: name,
	}
}

func (p *RegularPerson) Name() string {
	return p.name
}
