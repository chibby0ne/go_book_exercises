package tempconv

import (
	"flag"
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter2/tempconv"
)

type celsiusFlag struct{ tempconv.Celsius }

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "C", "°C":
		f.Celsius = tempconv.Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = tempconv.FToC(tempconv.Fahrenheit(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

// CelsiusFlag defines a Celsius Flag with the specified name, default value
// and usage and returns the address of the flag variable. The Flag argument
// must have a quantity and a unit, e.g: "100C"
func CelsiusFlag(name string, value tempconv.Celsius, usage string) *tempconv.Celsius {
	f := celsiusFlag{value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}
