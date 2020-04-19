// Extend the jpeg program so that it converts any supported input format to
// any output format, using image.Decode to detect the input formatter and a
// flag to select the output format

package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"  // register GIF decoder
	"image/jpeg" // register JPEG decoder
	"image/png"  // register PNG decoder
	"io"
	"log"
	"os"
	"strings"
)

var formats = []string{"jpeg", "gif", "png"}

func main() {
	outformat := flag.String("outformat", "jpeg", fmt.Sprintf("Specifies the output format. Possible values: %v", formats))
	flag.Parse()
	if err := validateOutFormat(outformat); err != nil {
		log.Fatal(err)
	}
	if err := toJPEG(os.Stdin, os.Stdout, outformat); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
	}
}

func validateOutFormat(outformat *string) error {
	for _, format := range formats {
		if strings.ToLower(*outformat) == format {
			*outformat = strings.ToLower(*outformat)
			return nil
		}
	}
	return fmt.Errorf("invalid format specified: %v", *outformat)
}

func toJPEG(in io.Reader, out io.Writer, outformat *string) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "input formater =", kind)
	switch *outformat {
	case "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "gif":
		return gif.Encode(out, img, &gif.Options{NumColors: 256})
	case "png":
		return png.Encode(out, img)
	}
	return nil
}
