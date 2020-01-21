package main

import (
	"io"
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
	n, err = io.Copy(f, resp.Body)
	// on many file systems, notably NFS write errors are not reported
	// immediatelly but may be postpoed until the file is closed.  Failure to
	// check the result of the close operation could cause serious data loss to
	// go unnoticed
	// That's why we are not deferring f.Close().
	// However if both io.Copy and
	// f.Close() fail we should prefer to report the error from io.Copy since
	// it occurred first and is more likely to tell us the root cause
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return local, n, err
}

func main() {
	for _, v := range os.Args[1:] {
		fetch(v)
	}
}
