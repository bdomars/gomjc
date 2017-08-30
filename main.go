package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bdomars/gomjc/microjava"
)

func main() {
	fmt.Println("gomjc v0.1a")

	if len(os.Args) < 2 {
		fmt.Println("err: no input specified")
		os.Exit(1)
	}

	infile, err := os.Open(os.Args[1])
	defer infile.Close()
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(infile)
	scanner := microjava.NewScanner(reader)
	for {
		token := scanner.NextToken()
		fmt.Println(token)
		if token.ErrorMsg != "" {
			panic(token.ErrorMsg)
		}
		if token.IsEOF() {
			break
		}
	}
}
