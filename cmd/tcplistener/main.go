package main

import (
	"fmt"
	"github.com/Ell534/tcpToHttp/internal/request"
	"log"
	"net"
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

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error when calling RequestFromReader: %s\n", err.Error())
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
	}
}
