package main

import (
	"github.com/adonovan/gopl.io/ch8/thumbnail"
)

// This function has a subtle bug.  When it encounters the first non-nil error,
// it returns the error to the caller, leaving no goroutine draining the errors
// channel. Each remaining worker goroutine will block forever when it tries
// tot send a value on that channel, and will never terminate.
func makeThumbnails4(filenames []string) error {
	errors := make(chan error)

	for _, f := range filenames {
		go func(f string) {
			_, err := thumbnail.ImageFile(f)
			errors <- err
		}(f)
	}
	for range filenames {
		if err := <-errors; err != nil {
			return err // NOTE: incorrect: gorutine leak!
		}
	}
	return nil
}

func main() {

}
