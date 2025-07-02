package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on port 42069")
	for {
		// wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("A connection has been accepted from", conn.RemoteAddr())

		linesChannel := getLinesChannel(conn)

		for line := range linesChannel {
			fmt.Println("read:", line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(connection net.Conn) <-chan string {
	lines := make(chan string)
	go func() {
		defer connection.Close()
		defer close(lines)

		var currentLine string
		for {
			bytes := make([]byte, 8, 8)
			n, err := connection.Read(bytes)
			if err != nil {
				if currentLine != "" {
					lines <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error reading from file: %s\n", err.Error())
				return
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
