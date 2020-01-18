package main

import (
	"bufio"
	"fmt"
	"os"
)

var InFile *bufio.Reader

func main() {
	if len(os.Args) != 2 {
		fatal(fmt.Sprintf("usage: %s infile\n", os.Args[0]))
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fatal("unable to open file %s: %v\n", os.Args[1], err)
	}
	InFile = bufio.NewReader(file)
	scanfile()
}

var tokenStrings = []string{"+", "-", "*", "/", "intlit"}

func scanfile() {
	t := &Token{}
	for scan(t) {
		fmt.Printf("Token %s", tokenStrings[t.token])
		if t.token == TokenIntLiteral {
			fmt.Printf(", value %d", t.value)
		}
		fmt.Printf("\n")
	}
}

func fatal(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(1)
}
