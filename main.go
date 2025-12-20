package main

import (
	"fmt"
	"io"
	"os"
) 

func main () {
	f, err := os.Open("messages.txt")
	if err != nil {
		panic(err) 
	}

	b := make([]byte, 8)
	for {
		n, err := f.Read(b)
		fmt.Printf("read: %s\n", b[:n]) 
		
		// error handling
		if err == io.EOF { // end of file stream
			break
		} else if err != nil { // other potential errors
			fmt.Printf("Error reading %v\n", err)
			break 
		}
	}
} 
