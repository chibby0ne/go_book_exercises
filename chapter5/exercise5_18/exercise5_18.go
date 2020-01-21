// Without changing its behaviour rewrite the fetch function to use defer to close the writable file.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	n, err = io.Copy(f, resp.Body)
	// on many file systems, notably NFS write errors are not reported
	// immediatelly but may be postpoed until the file is closed.  Failure to
	// check the result of the close operation could cause serious data loss to
	// go unnoticed
	// That's why we are not deferring f.Close().
	// However if both io.Copy and
	// f.Close() fail we should prefer to report the error from io.Copy since
	// it occurred first and is more likely to tell us the root cause
	return local, n, err
}

func main() {
	for _, url := range os.Args[1:] {
		filename, n, err := fetch(url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("url: %v, filename: %v, size: %v bytes\n", url, filename, n)
	}
}
