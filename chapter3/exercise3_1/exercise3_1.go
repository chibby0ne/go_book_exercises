// If the function f returns a non-finite float64 value, the SVG file will contain invalid <polygon> elements (althought many SVG renderers handle this gracefully). Modify the program to skip invalid polygons
package main

import (
	"fmt"
	"math"
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
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			var ax, ay, bx, by, cx, cy, dx, dy float64
			var ok bool
			if ax, ay, ok = corner(i+1, j); !ok {
				continue
			}
			if bx, by, ok = corner(i, j); !ok {
				continue
			}
			if cx, cy, ok = corner(i, j+1); !ok {
				continue
			}
			if dx, dy, ok = corner(i+1, j+1); !ok {
				continue
			}
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n", ax, ay, bx, by, cx, cy, dx, dy)
		}

	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64, bool) {
	// Find point (x, y) at corner of cell (i, j)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z
	z, ok := f(x, y)
	if !ok {
		return 0, 0, false
	}

	// Porject (x, y, z) isometrically onto 2-D SVG canvas (sx, sy)
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, true
}

func f(x, y float64) (float64, bool) {
	ok := true
	r := math.Hypot(x, y) // distance from (0, 0)
	ret := math.Sin(r) / r
	if math.IsInf(ret, 0) || math.IsNaN(ret) {
		ok = false
	}
	return ret, ok
}
