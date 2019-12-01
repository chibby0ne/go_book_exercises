// Another simple fractal uses Newton's method to find complex solution to a
// function such as z^4 - 1. Shade each starting point by the number of
// iterations required to get close to one of the four roots. Color each point
// by the root it approaches
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"math/rand"
	"os"
	"time"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
	iterations             = 200
	epsilon                = 1e-7
)

func main() {
	rand.Seed(time.Now().Unix())
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z
			img.Set(px, py, z4Minus1(z))
		}
	}
	png.Encode(os.Stdout, img)
}

func z4Minus1(z complex128) color.Color {
	f := func(p complex128) complex128 {
		return p*p*p*p - 1
	}
	fprime := func(p complex128) complex128 {
		return 4 * p * p * p
	}
	return newton(z, f, fprime)
}

var rootColor = map[complex128]color.RGBA{}

var colors = []color.RGBA{
	{44, 130, 201, 255},  // blue mariner
	{240, 52, 52, 255},   // red pomegranade
	{42, 187, 155, 255},  // green niagara
	{253, 227, 167, 255}, // orange cape honey
}

func newton(z0 complex128, f, fprime func(complex128) complex128) color.Color {
	z := z0
	for i := 0; i < iterations; i++ {
		if cmplx.Abs(f(z)) < epsilon {
			root := complex(round(real(z), 4), round(imag(z), 4))
			log.Printf("root is %v\n", root)
			color, ok := rootColor[root]
			if !ok {
				color = updateRootColors(root)
			}
			return chooseColor(i, color)
		}
		if fprime(z) == 0 {
			log.Printf("No solution found\n")
			return color.Black
		}
		z -= f(z) / fprime(z)
	}
	return color.Black
}

func updateRootColors(root complex128) color.RGBA {
	length := len(colors)
	if length == 0 {
		panic("no more colors left. Too many roots, due to too small epsilon or too small rounding?")
	}
	index := rand.Intn(length)
	rootColor[root] = colors[index]
	if index == length-1 {
		colors = colors[:length-1]
	} else if index == 0 {
		colors = colors[1:]
	} else {
		copy(colors[index:], colors[index+1:])
		colors = colors[:length-1]
	}
	return rootColor[root]
}

func chooseColor(i int, c color.RGBA) color.Color {
	y, cb, cr := color.RGBToYCbCr(c.R, c.G, c.B)
	// Need to use Log to make the difference in color more remarkable
	ratio := math.Log(float64(i)) / math.Log(iterations)
	y -= uint8(float64(y) * ratio)
	return color.YCbCr{y, cb, cr}
}

// round takes a float and truncates it to the number of digits but also rounds
// it by adding 5 to the decimal place after the digits position
func round(f float64, digits int) float64 {
	if math.Abs(f) < 0.5 {
		return 0
	}
	divisor := math.Pow10(digits)
	new_f := math.Trunc(f*divisor+math.Copysign(0.5, f)) / divisor
	return new_f
}
