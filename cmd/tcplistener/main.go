package main

import (
	"fmt"
	"log"
	"net"

	"github.com/kalim-Asim/http-server/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		// wait for connection
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		rl := req.RequestLine
		fmt.Println("Request line:")
		fmt.Printf(" - Method: %s\n", rl.Method)
		fmt.Printf(" - Target: %s\n", rl.RequestTarget)
		fmt.Printf(" - Version: %s\n", rl.HttpVersion)

		headers := req.Headers
		fmt.Println("Headers:")
		headers.PrintHeaders()


		fmt.Println("Body:")
		fmt.Printf("%s\n", req.Body)
	}
}
