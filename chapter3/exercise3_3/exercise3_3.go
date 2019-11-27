// Color each polygon based on its height, so that the peaks are colored red (#ff0000) and the valleys blue (#0000ff)
package main

import (
	"fmt"
	"math"
	"os"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange,..,xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	zmin, zmax := findMinMax()
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			if isNotValid(ax, ay) || isNotValid(bx, by) || isNotValid(cx, cy) || isNotValid(dx, dy) {
				continue
			}
			fmt.Printf("<polygon style='stroke: %s; fill: #222222' points='%g,%g %g,%g %g,%g %g,%g'/>\n", color(i, j, zmin, zmax), ax, ay, bx, by, cx, cy, dx, dy)
		}

	}
	fmt.Println("</svg>")
}

func findMinMaxCell(row, column int) (min, max float64) {
	min = math.NaN()
	max = math.NaN()
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			x := xyrange * (float64(i+row)/cells - 0.5)
			y := xyrange * (float64(j+column)/cells - 0.5)
			z := f(x, y)
			if math.IsNaN(min) || z < min {
				min = z
			}
			if math.IsNaN(max) || z > max {
				max = z
			}
		}
	}
	return min, max
}

func findMinMax() (min, max float64) {
	min = math.NaN()
	max = math.NaN()
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			temp_min, temp_max := findMinMaxCell(i, j)
			if math.IsNaN(min) || temp_min < min {
				min = temp_min
			}
			if math.IsNaN(max) || temp_max > max {
				max = temp_max
			}
		}
	}
	return min, max
}

func color(i, j int, zmin, zmax float64) string {
	min, max := findMinMaxCell(i, j)
	var color string
	if math.Abs(max) > math.Abs(min) {
		red := math.Exp(math.Abs(max)) / math.Exp(math.Abs(zmax)) * 255
		if red > 255 {
			red = 255
		}
		color = fmt.Sprintf("#%02x0000", int(red))
	} else {
		blue := math.Exp(math.Abs(min)) / math.Exp(math.Abs(zmin)) * 255
		if blue > 255 {
			blue = 255
		}
		color = fmt.Sprintf("#0000%02x", int(blue))
	}
	return color
}

func isNotValid(x, y float64) bool {
	return math.IsInf(x, 0) || math.IsNaN(x) || math.IsInf(y, 0) || math.IsNaN(y)
}

func corner(i, j int) (float64, float64) {
	// Find point (x, y) at corner of cell (i, j)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z
	z := f(x, y)

	// Porject (x, y, z) isometrically onto 2-D SVG canvas (sx, sy)
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0, 0)
	ret := math.Sin(r) / r
	return ret
}
