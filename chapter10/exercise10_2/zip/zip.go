package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	genarchive "github.com/chibby0ne/go_book_exercises/chapter10/exercise10_2"
)

type Reader struct {
	reader *zip.Reader
}

func (r *Reader) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	var read int
	for _, file := range r.reader.File {
		if len(b) == 0 {
			break
		}
		f, err := file.Open()
		if err != nil {
			return 0, fmt.Errorf("zip Read: %s", err)
		}
		defer f.Close()
		n, err := f.Read(b)
		read += n
		if err != nil && err != io.EOF {
			return read, err
		}
		b = b[n:]
	}
	return read, nil
}

func init() {
	genarchive.RegisterFormat("zip", 0, "PK", NewReader)
}

func NewReader(f *os.File) (io.Reader, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("new zip reader: %s", err)
	}
	r, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return nil, fmt.Errorf("new zip reader: %s", err)
	}
	return &Reader{r}, nil
}
