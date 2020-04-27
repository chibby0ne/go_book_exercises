package genarchive_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	genarchive "github.com/chibby0ne/go_book_exercises/chapter10/exercise10_2"
	_ "github.com/chibby0ne/go_book_exercises/chapter10/exercise10_2/tar"
	_ "github.com/chibby0ne/go_book_exercises/chapter10/exercise10_2/zip"
)

func TestBlankImport(t *testing.T) {
	inputs := []string{"example.tar", "example.zip"}
	for _, input := range inputs {
		f, err := os.Open(filepath.Join("testdata", input))
		if err != nil {
			t.Error(err)
		}
		defer f.Close()
		log.Println(f.Name())
		r, err := genarchive.Open(f)
		if err != nil {
			t.Error(err)
		}
		var bb []byte
		bb = make([]byte, 100)
		n, err := r.Read(bb)
		log.Printf("bb read %v bytes: %v", n, string(bb))
		if err != nil {
			t.Error(err)
		}
		// The order of the files read depends on the order of the files that
		// were given the when the zip/tar were created, perhaps a more robust
		// test is to check against the size of the read bytes. In my case I
		// simply created it in the order I expect them:
		//
		// tar cvf example.tar dir1 file2.txt
		// zip -0 example.zip dir1/file1.txt file2.txt
		//
		// where file1.txt contains: 'hello' and file2.txt contains: 'world'
		//
		expected := "helloworld"
		if string(bb[:n]) != expected {
			t.Error("not equal!")
		}
	}
}
