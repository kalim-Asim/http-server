package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	f, err := os.Open("messages.txt")
	if err != nil {
		panic(err) 
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
} 
