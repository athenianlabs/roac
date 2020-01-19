package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	InFile  *bufio.Reader
	OutFile *bufio.Writer
)

func main() {
	if len(os.Args) != 2 {
		fatal(fmt.Sprintf("usage: %s infile\n", os.Args[0]))
	}
	inFile, err := os.Open(os.Args[1])
	if err != nil {
		fatal("unable to open file %s: %v\n", os.Args[1], err)
	}
	defer inFile.Close()
	InFile = bufio.NewReader(inFile)

	outFile, err := os.Create("out.s")
	if err != nil {
		fatal("unable to create out.s: %v\n", err)
	}
	defer outFile.Close()
	OutFile = bufio.NewWriter(outFile)

	scan(CurrentToken) // Get the first token from the input
	n := binexpr(0)    // Parse the expression in the file
	fmt.Printf("%d\n", interpretAST(n))
	generatecode(n)
	if err := OutFile.Flush(); err != nil {
		fatal("unable to write to out.s: %v\n", err)
	}
}

func fatal(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
	os.Exit(1)
}
