package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("failed to open file %s due to error: %s\n", inputFilePath, err)
	}

	fmt.Printf("Currently reading data from: %s\n", inputFilePath)
	fmt.Println("=================================")

	linesChannel := getLinesChannel(file)

	for line := range linesChannel {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)

		var currentLine string
		for {
			bytes := make([]byte, 8, 8)
			n, err := f.Read(bytes)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error reading from file: %s\n", err.Error())
				break
			}

			stringToPrint := string(bytes[:n])
			parts := strings.Split(stringToPrint, "\n")

			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}

			currentLine += parts[len(parts)-1]
		}
	}()
	return lines
}
