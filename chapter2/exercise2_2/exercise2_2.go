package main

import (
	"fmt"
	"github.com/chibby0ne/go_book_exercises/chapter2/tempconv"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Meter float64
type Feet float64

func MToFeet(m Meter) Feet {
	return Feet(m * 3.2804)
}

func FeetToM(f Feet) Meter {
	return Meter(f / 3.2804)
}

func usage(program string) string {
	return fmt.Sprintf(`usage: %s value unit`, program)
}

var exp = regexp.MustCompile(`(?P<value>[0-9]+\.?[0-9]*)(?P<unit>[m|M|Meter|meter|"|feet|Feet|f|F|Fahrenheit|fahrenheit|c|C|Celsius|celsius])`)

func main() {
	var input string
	if len(os.Args) == 2 {
		input = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Print(usage(os.Args[0]))
		os.Exit(1)
	} else {
		fmt.Println("Enter value and unit with no space: eg: 12m")
		fmt.Scanf("%s", &input)
	}
	match := exp.FindStringSubmatch(input)
	result := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	value, err := strconv.ParseFloat(result["value"], 64)
	unit := strings.ToLower(result["unit"])
	if err != nil {
		fmt.Printf("Error converting the value %v", err)
	}
	switch unit {
	case "m", "meter":
		distance_meter := Meter(value)
		fmt.Printf("%v m = %v\"\n", distance_meter, MToFeet(distance_meter))
	case "\"", "feet":
		distance_feet := Feet(value)
		fmt.Printf("%v\" = %v m\n", distance_feet, FeetToM(distance_feet))
	case "c", "celsius":
		temp_celsius := tempconv.Celsius(value)
		fmt.Printf("%v = %v\n", temp_celsius, tempconv.CToF(temp_celsius))
	case "f", "fahrenheit":
		temp_fahrenheit := tempconv.Fahrenheit(value)
		fmt.Printf("%v = %v\n", temp_fahrenheit, tempconv.FToC(temp_fahrenheit))
	default:
		fmt.Fprintf(os.Stderr, "No such case: %v\n", unit)
	}

}
