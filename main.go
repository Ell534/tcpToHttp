package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const inputFile = "messages.txt"

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("error opening file %s: %v\n", inputFile, err)
	}
	defer file.Close()

	fmt.Printf("reading data from input file: %s\n", inputFile)
	fmt.Println("======================================")

	for {
		bytes := make([]byte, 8)
		n, err := file.Read(bytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error reading file: %s\n", err.Error())
			break
		}
		strToPrint := string(bytes[:n])
		fmt.Printf("read: %s\n", strToPrint)
	}
}
