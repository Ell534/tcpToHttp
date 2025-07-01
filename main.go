package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("failed to open file %s due to error: %s\n", inputFilePath, err)
	}
	defer file.Close()

	fmt.Printf("Currently reading data from: %s\n", inputFilePath)
	fmt.Println("=================================")

	for {
		bytes := make([]byte, 8, 8)
		n, err := file.Read(bytes)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error reading from file: %s\n", err.Error())
			break
		}
		stringToPrint := string(bytes[:n])
		fmt.Printf("read: %s\n", stringToPrint)
	}
}
