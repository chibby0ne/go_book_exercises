// Write a web server that renders fractals and writes the image data to the client.
// Allow the client to specify the x, y, and zoom values as parameters to the HTTP request.

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"os"
	"strconv"
)

const (
	Xmin, Ymin, Xmax, Ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
	Iterations             = 200
	Contrast               = 15
	Zoom                   = 1
	Url                    = "localhost:8000"
	Usage                  = `Usage: query for something all of the keys "xmin", "ymin", "xmax", "ymax", "zoom"
Example: /?xmin=1&ymin=1&xmax=3,ymax=3,zoom=3
`
)

var Params = map[string]float64{
	"xmin": Xmin,
	"ymin": Ymin,
	"xmax": Xmax,
	"ymax": Ymax,
	"zoom": Zoom,
}

func main() {
	flag.Usage = func() {
		fmt.Println(Usage)
	}
	fmt.Printf("Listening to %v\n", Url)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(Url, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stdout, "%s %s %s\n", r.Method, r.URL, r.Proto)
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing the form: %v\n", err), http.StatusBadRequest)
		flag.Usage()
		return
	}
	for k, v := range Params {
		val := r.Form.Get(k)
		if val == "" {
			fmt.Fprintf(os.Stdout, "Setting %q to default: %v\n", k, v)
			Params[k] = v
			continue
		}
		res, err := strconv.ParseFloat(val, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing %q=%q to float64: %v", k, val, err), http.StatusBadRequest)
			flag.Usage()
			return
		}
		fmt.Fprintf(os.Stdout, "Setting %q = to %v\n", k, res)
		Params[k] = res
	}

	if Params["xmax"] <= Params["xmin"] || Params["ymax"] <= Params["ymin"] {
		flag.Usage()
		http.Error(w, "xmax/ymax should be bigger than xmin/ymax respectively", http.StatusBadRequest)
		return
	}

	xlength := Params["xmax"] - Params["xmin"]
	ylength := Params["ymax"] - Params["ymin"]

	xmid := Params["xmin"] + xlength/2
	ymid := Params["ymin"] + ylength/2

	Params["xmin"] = xmid - (xlength/2)/Params["zoom"]
	Params["ymin"] = ymid - (ylength/2)/Params["zoom"]

	Params["xmax"] = xmid + (xlength/2)/Params["zoom"]
	Params["ymax"] = ymid + (ylength/2)/Params["zoom"]
	if err := makeMandelbrot(w); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering the png image: %v", err), http.StatusInternalServerError)
	}
}

func makeMandelbrot(output io.Writer) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(Params["ymax"]-Params["ymin"]) + Ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(Params["xmax"]-Params["xmin"]) + Xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z
			img.Set(px, py, mandelbrotComplex128(z))
		}
	}
	return png.Encode(output, img)
}

func mandelbrotComplex128(z complex128) color.Color {
	var v complex128
	for n := uint8(0); n < Iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - Contrast*n}
		}
	}
	return color.Black
}
