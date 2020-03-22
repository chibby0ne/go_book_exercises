package flag

import (
	"flag"
)

const (
	DefaultControlPort int64 = 21
	DefaultDataPort    int64 = 20
)

var (
	Port     = flag.Int64("port", DefaultControlPort, "control port to listen to connections")
	DataPort = flag.Int64("dport", DefaultDataPort, "data port to create data connection from")
	Verbose  = flag.Bool("v", false, "verbose output")
)
