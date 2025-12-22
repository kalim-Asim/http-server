package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
) 

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		
		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			
			// read 8 bytes but print a line
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i == -1 {
				str += string(data)
			} else {
				str += string(data[:i])
				out <- str
				str = string(data[i+1:])
			}

			// error handling
			if err == io.EOF { 
				break
			} else if err != nil { 
				fmt.Printf("Error reading %v\n", err)
				break 
			}
		}

		if str != "" {
			out <- str
			str = ""
		}
	}()

	return out 
}

func main () {
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

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}
} 
