package exercise13_4

import (
	"bytes"
	"compress/bzip2"
	"io"
	"testing"
)

const (
	DataToCompress = "asdfasdf"
	Iterations     = 1e4
)

func TestWriter(t *testing.T) {
	compressed, originalUncompressed := &bytes.Buffer{}, &bytes.Buffer{}
	writer, err := NewWriter(compressed)
	if err != nil {
		t.Errorf("error creating newWriter: %v", err)
	}

	// Write the same thing to an original uncompressed buffer and a buffer that is to be compressed
	for i := 0; i < Iterations; i++ {
		writer.Write([]byte(DataToCompress))
		originalUncompressed.Write([]byte(DataToCompress))
	}
	err = writer.Close()
	if err != nil {
		t.Errorf("error closing the writer: %v", err)
	}

	// Uncompresses the compressed bytes
	reader := bzip2.NewReader(compressed)

	uncompressed := &bytes.Buffer{}
	_, err = io.Copy(uncompressed, reader)

	if err != nil {
		t.Errorf("error copying the reader to the buffer: %v", err)
	}

	if !bytes.Equal(uncompressed.Bytes(), originalUncompressed.Bytes()) {
		t.Errorf("bytes uncompressed after being compressed are not equal to the original bytes uncompressed")
	}

}
