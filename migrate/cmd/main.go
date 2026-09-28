package main

import (
	"flag"
	"fmt"
	"os"

	"migrate/runner"
)

func main() {
	from := flag.String("from", "", "Starting version")
	to := flag.String("to", "", "Target version")
	flag.Parse()

	if *from == "" || *to == "" || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: migrate -from <version> -to <version>")
		os.Exit(2)
	}

	if err := runner.Run(*from, *to); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
