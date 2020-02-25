package main

import (
	"bytes"
	"io"
	"os"
	"time"
)

func main() {
	var w io.Writer
	w = os.Stdout
	w = new(bytes.Buffer)
	w = time.Second

	var rwc io.ReadWriteCloser
	rwc = os.Stdout
	rwc = new(bytes.Buffer)

	w = rwc
	rwc = w

}
