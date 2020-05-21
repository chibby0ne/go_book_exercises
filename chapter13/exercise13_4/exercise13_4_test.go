package exercise13_4

import (
	"bytes"
	"compress/bzip2"
	"io"
	"testing"
)

const DataToCompress = "asdfasdf"

func TestWriter(t *testing.T) {
	compressed, originalUncompressed := &bytes.Buffer{}, &bytes.Buffer{}
	writer, err := NewWriter(compressed)
	if err != nil {
		t.Errorf("error creating newWriter: %v", err)
	}

	for i := 0; i < 1e4; i++ {
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
		t.Errorf("value uncompressed is not equal to the original value uncomprossed")
	}

}
