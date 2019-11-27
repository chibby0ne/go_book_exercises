// Implement full color Mandlebrot set using the fucntion image.NEwRGBA and the ttype color.RGBA or color.YCbCr
package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z
			var subpixels [4]color.Color
			for i := 0; i < 2; i++ {
				for j := 0; j < 2; j++ {
					subpixels[i*2+j] = mandelbrot(z)
				}
			}
			img.Set(px, py, average(subpixels))
		}
	}
	png.Encode(os.Stdout, img)
}

func average(subpixels [4]color.Color) color.Color {
	var r, g, b, a uint16
	for i := 0; i < 4; i++ {
		r_, g_, b_, a_ := subpixels[i].RGBA()
		r += uint16(r_ / 4)
		g += uint16(g_ / 4)
		b += uint16(b_ / 4)
		a += uint16(a_ / 4)
	}
	return color.RGBA64{r, g, b, a}
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const constrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - constrast*n}
		}
	}
	return color.Black
}
