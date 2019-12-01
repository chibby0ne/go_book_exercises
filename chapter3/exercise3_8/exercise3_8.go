// Rendering fractals at high zoom level s demands gre at arithmetic pre cision.
// Implement the same fractal using four different representations of numbers: complex64, com-
// plex128, big.Float, and big.Rat. (The latter two types are found in the math/big package.
// Float uses arbitrary but bounded-precision floating-point; Rat uses unbounded-precision
// rational numbers.) How do they compare in performance and memory usage? At what zoom
// levels do rendering artifacts become visible?

package main

import (
	// "flag"
	"image"
	"image/color"
	"image/png"
	// "log"
	"math/big"
	"math/cmplx"
	"os"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
	ITERATIONS             = 200
	CONTRAST               = 15
)

type complexBigFloat struct {
	real_ big.Float
	imag_ big.Float
}

func main() {

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z
			img.Set(px, py, mandelbrotBigFloat(z))
		}
	}
	png.Encode(os.Stdout, img)
}

func mandelbrotComplex64(z complex64) color.Color {
	var v complex64
	for n := uint8(0); n < ITERATIONS; n++ {
		v = v*v + z
		if cmplx.Abs(complex128(v)) > 2 {
			return color.Gray{255 - CONTRAST*n}
		}
	}
	return color.Black
}

func mandelbrotComplex128(z complex128) color.Color {
	var v complex128
	for n := uint8(0); n < ITERATIONS; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - CONTRAST*n}
		}
	}
	return color.Black
}

func multiplyComplexBigFloat(v, z complexBigFloat) complexBigFloat {
	// v = (a + bi), z = (c + di)
	// (a + bi)(c + di) = a * c + a * di + bi * c + b * d * i^2 = (ac - bd) + (ad + bc)i
	var real_part big.Float
	var imag_part big.Float
	real_part.Sub((&big.Float{}).Mul(&v.real_, &z.real_), (&big.Float{}).Mul(&v.imag_, &z.imag_))
	imag_part.Add((&big.Float{}).Mul(&v.real_, &z.imag_), (&big.Float{}).Mul(&v.imag_, &z.real_))
	return complexBigFloat{real_: real_part, imag_: imag_part}
}

func addComplexBigFloat(a, b complexBigFloat) complexBigFloat {
	return complexBigFloat{real_: *(&big.Float{}).Add(&a.real_, &b.real_), imag_: *(&big.Float{}).Add(&a.imag_, &b.imag_)}
}

func absoluteValueComplexBigFloat(z complexBigFloat) *big.Float {
	var result big.Float
	result.Add((&big.Float{}).Mul(&z.real_, &z.real_), (&big.Float{}).Mul(&z.imag_, &z.imag_))
	return result.Sqrt(&result)
}

func mandelbrotBigFloat(z complex128) color.Color {
	var z_big_float = complexBigFloat{real_: *big.NewFloat(real(z)), imag_: *big.NewFloat(imag(z))}
	var v complexBigFloat
	var two_big_float = big.NewFloat(2)
	var temp_v complexBigFloat
	for n := uint8(0); n < ITERATIONS; n++ {
		if n != 0 {
			temp_v = multiplyComplexBigFloat(v, v)
			v = addComplexBigFloat(temp_v, z_big_float)
		} else {
			v = z_big_float
		}
		if absoluteValueComplexBigFloat(v).Cmp(two_big_float) == 1 {
			return color.Gray{255 - CONTRAST*n}
		}
	}
	return color.Black
}

type complexBigRat struct {
	real_ big.Rat
	imag_ big.Rat
}

func multiplyComplexBigRat(v, z complexBigRat) complexBigRat {
	// v = (a + bi), z = (c + di)
	// (a + bi)(c + di) = a * c + a * di + bi * c + b * d * i^2 = (ac - bd) + (ad + bc)i
	var real_part big.Rat
	var imag_part big.Rat
	real_part.Sub((&big.Rat{}).Mul(&v.real_, &z.real_), (&big.Rat{}).Mul(&v.imag_, &z.imag_))
	imag_part.Add((&big.Rat{}).Mul(&v.real_, &z.imag_), (&big.Rat{}).Mul(&v.imag_, &z.real_))
	return complexBigRat{real_: real_part, imag_: imag_part}
}

func addComplexBigRat(a, b complexBigRat) complexBigRat {
	return complexBigRat{real_: *(&big.Rat{}).Add(&a.real_, &b.real_), imag_: *(&big.Rat{}).Add(&a.imag_, &b.imag_)}
}

func absoluteValueComplexBigRat(z complexBigRat) *big.Rat {
	var result big.Rat
	result.Add((&big.Rat{}).Mul(&z.real_, &z.real_), (&big.Rat{}).Mul(&z.imag_, &z.imag_))
	return &result
}

func mandelbrotBigRat(z complex128) color.Color {
	var z_big_rat = complexBigRat{real_: *(new(big.Rat).SetFloat64(real(z))), imag_: *(new(big.Rat).SetFloat64(imag(z)))}
	var v complexBigRat
	var four_big_rat = big.NewRat(4, 1)
	var temp_v complexBigRat
	for n := uint8(0); n < ITERATIONS; n++ {
		if n != 0 {
			temp_v = multiplyComplexBigRat(v, v)
			v = addComplexBigRat(temp_v, z_big_rat)
		} else {
			v = z_big_rat
		}
		if absoluteValueComplexBigRat(v).Cmp(four_big_rat) == 1 {
			return color.Gray{255 - CONTRAST*n}
		}
	}
	return color.Black
}
