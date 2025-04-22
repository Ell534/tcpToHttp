package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	var currentLine string

	for {
		bytes := make([]byte, 8)
		n, err := file.Read(bytes)

		if err != nil {
			if currentLine != "" {
				fmt.Printf("read: %s\n", currentLine)
				currentLine = ""
			}
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error reading file: %s\n", err.Error())
			break
		}

		str := string(bytes[:n])
		parts := strings.Split(str, "\n")

		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", currentLine, parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}
}
