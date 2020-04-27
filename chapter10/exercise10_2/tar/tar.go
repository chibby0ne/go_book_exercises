package tar

import (
	"archive/tar"
	"io"
	"log"
	"os"

	genarchive "github.com/chibby0ne/go_book_exercises/chapter10/exercise10_2"
)

type Reader struct {
	reader *tar.Reader
}

// According to the definition of (*tar.Reader).Read, trying to read on one of the special types returns (0, io.EOF)
func isSpecialType(header *tar.Header) bool {
	return header.Typeflag == tar.TypeLink || header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeChar || header.Typeflag == tar.TypeBlock || header.Typeflag == tar.TypeDir || header.Typeflag == tar.TypeFifo
}

func (r *Reader) Read(b []byte) (int, error) {
	var bytesRead int
	for len(b) > 0 {
		header, err := r.reader.Next()
		if err == io.EOF {
			log.Println("Breaking from err == io.EOF in first next")
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// Keep looking for the next regular file
		for isSpecialType(header) {
			header, err = r.reader.Next()
			if err == io.EOF {
				log.Println("Breaking from err == io.EOF in first next")
				break
			}
			if err != nil {
				log.Fatal(err)
			}
		}
		num, err := r.reader.Read(b)
		log.Printf("after reading, num: %d, err: %s, b: %v", num, err, b)
		bytesRead += num
		log.Printf("bytesRead: %v", bytesRead)
		// Read the whole file? Log it
		if err == io.EOF {
			log.Printf("Read the whole file: %v", header.Name)
		}
		// There were some errors and it wasn't EOFJk
		if err != nil && err != io.EOF {
			return bytesRead, err
		}
		b = b[num:]
	}
	return bytesRead, nil
}

func init() {
	genarchive.RegisterFormat("tar", 257, "ustar", NewReader)
}

func NewReader(f *os.File) (io.Reader, error) {
	// log.Printf("new reader file: %p", f)
	// sanityCheck(f)
	// log.Printf("file: %v", f.Name())
	return &Reader{tar.NewReader(f)}, nil
}

func sanityCheck(f *os.File) {
	log.Printf("sanityCheck: file: %p", f)
	log.Printf("sanity check: file: %v", f.Name())
	reader := tar.NewReader(f)
	b := make([]byte, 7)
	log.Printf("sanity check: len of b: %v, cap of b: %v\n", len(b), cap(b))
outer:
	for {
		header, err := reader.Next()
		if err == io.EOF {
			log.Println("sanity check: Breaking from err == io.EOF in first next")
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for isSpecialType(header) {
			header, err = reader.Next()
			if err == io.EOF {
				log.Println("sanity check: Breaking from err == io.EOF in second next")
				break outer
			}
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("sanity check: header: %+v\n", header)
		n, err := reader.Read(b)
		if err == io.EOF {
			log.Println(err)
		}
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		log.Printf("n: %v", n)
		log.Printf(string(b))
		b = b[n:]
	}
}
