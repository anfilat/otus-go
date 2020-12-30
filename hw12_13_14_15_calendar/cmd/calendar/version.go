package main

import (
	"flag"
	"fmt"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func printVersion() {
	fmt.Printf("Calendar %s release (%s) built on %s\n", release, gitHash, buildDate)
}

func isVersionCommand() bool {
	for _, name := range flag.Args() {
		if name == "version" {
			return true
		}
	}
	return false
}
