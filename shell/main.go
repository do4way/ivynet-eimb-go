package main

import (
	"flag"
	"fmt"
	"os"
)

const closeFd = ^uintptr(0)

var (
	help       = flag.Bool("help", false, "show usage help and quit")
	logpath    = flag.String("log", "", "a file to write debug log to")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	cmd        = flag.Bool("c", false, "take first argument as a command to execute")
	forked     = flag.Int("forked", 0, "how many times the daemon has forked")
)

func usage() {
	fmt.Println("usage: esh [flags] [script]")
	fmt.Println("flags:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if *help {
		usage()
		os.Exit(0)
	}

	fmt.Printf("Not implemented yet, %v", args)
}
