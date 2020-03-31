package flag

import (
	"flag"
)

const (
	DefaultControlPort int64 = 21
)

var (
	Port    = flag.Int64("port", DefaultControlPort, "control port to listen to connections")
	Verbose = flag.Bool("v", false, "verbose output")
)
