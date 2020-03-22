package log

import (
	"log"

	"github.com/chibby0ne/go_book_exercises/chapter8/exercise8_2/flag"
)

func LogVerbose(a ...interface{}) {
	if *flag.Verbose {
		log.Print(a...)
	}
}

func LogfVerbose(format string, a ...interface{}) {
	if *flag.Verbose {
		log.Printf(format, a...)
	}
}

func Fatal(a ...interface{}) {
	log.Fatal(a...)
}

func Fatalf(format string, a ...interface{}) {
	log.Fatalf(format, a...)
}
