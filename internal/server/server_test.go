package server

import (
	"time"
	"net"
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
    port := 8080
    srv, err := Serve(port)
    if err != nil {
        t.Fatalf("Could not start server: %v", err)
    }
    defer srv.Close()

    // Give the server a moment to start the goroutine
    time.Sleep(100 * time.Millisecond)

    // Attempt to connect as a client
    conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        t.Fatalf("Could not connect to server: %v", err)
    }
    defer conn.Close()

    // Read the response
    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        t.Fatalf("Read error: %v", err)
    }
    
    fmt.Printf("Server responded: %s\n", string(buf[:n]))
}
