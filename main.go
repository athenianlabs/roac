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
	scan(CurrentToken) // Get the first token from the input
	n := binexpr()     // Parse the expression in the file
	fmt.Printf("%d\n", interpretAST(n))
}

var tokenStrings = []string{"+", "-", "*", "/", "intlit"}

func fatal(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(1)
}
