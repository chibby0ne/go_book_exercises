package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

var red = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
var green = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
var blue = color.RGBA{0x00, 0x00, 0xFF, 0xFF}
var yellow = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
var aquamarine = color.RGBA{0x00, 0xFF, 0xFF, 0xFF}
var pink = color.RGBA{0xFF, 0x00, 0x7B, 0xFF}

var palette = []color.Color{color.White, color.Black, red, green, blue, yellow, aquamarine, pink}

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas convers
		nframes = 64    // number of animation
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		colorIndex := uint8(rand.Int31n(int32(len(palette))))
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), colorIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
