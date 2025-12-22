package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"fmt"
)

func main() {
	address := "localhost:42069" // ip address and port 

	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(" > ")
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue 
		}

		if _, err := conn.Write([]byte(str)); err != nil {
			fmt.Printf("Error sending to udp message: %v", err)
			continue
		}
	}
}
