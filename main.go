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

	fmt.Printf("reading data from input file: %s\n", inputFile)
	fmt.Println("======================================")

	linesChan := getLinesChannel(file)

	for line := range linesChan {
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	result := make(chan string)

	go func() {
		defer f.Close()
		defer close(result)
		var currentLine string

		for {
			bytes := make([]byte, 8)
			n, err := f.Read(bytes)

			if err != nil {
				if currentLine != "" {
					result <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}

			str := string(bytes[:n])
			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				result <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()
	return result
}
