package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
) 

func main () {
	f, err := os.Open("messages.txt")
	if err != nil {
		panic(err) 
	}

	data := make([]byte, 8)
	str := ""
	for {
		n, err := f.Read(data)
		
		// read 8 bytes but print a line
		data = data[:n]
		if i := bytes.IndexByte(data, '\n'); i == -1 {
			str += string(data)
		} else {
			str += string(data[:i])
			fmt.Printf("read: %s\n", str)
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
		fmt.Printf("read: %s\n", str)
		str = ""
	}
} 
