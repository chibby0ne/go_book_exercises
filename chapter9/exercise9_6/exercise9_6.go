// Measure how the performance of a compute-bound parallel program varies with
// GOMAXPOCS. What is the optimal value on your computer? How many CPUs does
// your computer have?
// Optimal value is the number of logical cores

package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
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

type zFunc func(x, y float64) float64

func main() {
	var f zFunc = drop
	if len(os.Args) == 2 {
		switch strings.ToLower(os.Args[1]) {
		case "flatdrop":
			f = flatdrop
		case "saddle":
			f = saddle
		case "drop":
			f = drop
		default:
			fmt.Fprintf(os.Stderr, "usage: %s flatdrop | saddle | drop\n", os.Args[0])
			os.Exit(1)
		}
	}
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>\n", width, height)
	workers := runtime.GOMAXPROCS(runtime.NumCPU())
	var channels []chan string

	// Workers goroutines -> each one doing a section of the whole 100 x 100 cells
	for n := 0; n < workers; n++ {
		start, end := calculateRanges(n, workers)
		ch := make(chan string, (end-start)*cells)
		go func(start, end int, ch chan string) {
			for i := start; i < end; i++ {
				for j := 0; j < cells; j++ {
					ax, ay := corner(i+1, j, f)
					bx, by := corner(i, j, f)
					cx, cy := corner(i, j+1, f)
					dx, dy := corner(i+1, j+1, f)
					if isNotValid(ax, ay) || isNotValid(bx, by) || isNotValid(cx, cy) || isNotValid(dx, dy) {
						continue
					}
					ch <- fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n", ax, ay, bx, by, cx, cy, dx, dy)
				}
			}
			close(ch)
		}(start, end, ch)
		channels = append(channels, ch)
	}
	// Writer goroutine.
	// In svg the order of the polygon points elements matter to much of my frustration :(
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for _, ch := range channels {
			for str := range ch {
				fmt.Print(str)
			}
		}
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("</svg>")
}

func calculateRanges(n, workers int) (startRange int, endRange int) {
	startRange = n * (cells / workers)
	endRange = cells
	if n != workers-1 {
		endRange = (n + 1) * (cells / workers)
	}
	return
}

func flatdrop(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0, 0)
	ret := math.Sin(r) / xyrange
	return ret
}

func drop(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0, 0)
	ret := math.Sin(r) / r
	return ret
}

func saddle(x, y float64) float64 {
	ret := math.Pow(x, 2)/width + math.Pow(y, 2)/height
	return ret
}

func isNotValid(x, y float64) bool {
	return math.IsInf(x, 0) || math.IsNaN(x) || math.IsInf(y, 0) || math.IsNaN(y)
}

func corner(i, j int, f zFunc) (float64, float64) {
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
