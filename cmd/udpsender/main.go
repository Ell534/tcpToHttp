package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("error resolving UDP address: %s\n", err.Error())
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("error establishing udp connection: %s\n", err.Error())
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading line from reader: %s\n", err.Error())
		}
		_, err = udpConn.Write([]byte(line))
		if err != nil {
			log.Fatalf("error writing line to the udp connection: %s\n", err.Error())
		}
	}
}
