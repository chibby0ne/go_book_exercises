// Define a generic archive file-reading function capable of reading ZIP files
// (archive/zip) and POSIX tar files (archive/tar). Use a registration
// mechanism similar to the one described above so that support for each file
// format can be plugged in using blank imports.
package genarchive

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

type readerFunc func(*os.File) (io.Reader, error)

type format struct {
	name        string
	magicOffset int
	magicString string
	read        readerFunc
}

var (
	formatsMu sync.Mutex
	formats   []*format
	ErrFormat = errors.New("archive: unknown format")
)

func RegisterFormat(name string, magicOffset int, magicString string, read readerFunc) {
	formatsMu.Lock()
	formats = append(formats, &format{name, magicOffset, magicString, read})
	formatsMu.Unlock()
}

func detectArchive(f *os.File) *format {
	r := bufio.NewReader(f)
	var matchedFormat *format
	for _, form := range formats {
		//  Read up to the magic string without advancing the reader
		peekedBytes, err := r.Peek(form.magicOffset + len(form.magicString))
		if err != nil {
			continue
		}
		if string(peekedBytes[form.magicOffset:]) == form.magicString {
			log.Printf("file: %v is of format: %v", f.Name(), form.name)
			matchedFormat = form
			break
		}
	}
	// return nil if we found no match or the actual format if we found something
	return matchedFormat
}

func Open(file *os.File) (io.Reader, error) {
	detectedArchive := detectArchive(file)
	// if it wasn't a known file archive format, return an error
	if detectedArchive == nil {
		return nil, ErrFormat
	}
	// Rewind the file pointer from the detectArchive bufio.NewReader's peek
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	// return the io.Reader of the file format
	return detectedArchive.read(file)
}
