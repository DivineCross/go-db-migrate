package main

import (
	"flag"
	"fmt"
	"os"

	"migrate/runner"
)

func main() {
	from := flag.Int("from", -1, "Starting step sequence")
	to := flag.Int("to", -1, "Target step sequence")
	flag.Parse()

	if *from < 0 || *to < 0 || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: migrate -from <seq> -to <seq>")
		os.Exit(2)
	}

	if err := runner.Run(*from, *to); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
