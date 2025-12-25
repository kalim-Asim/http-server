package server

import (
	// "log"
	"net"
	"sync/atomic"
	// "net/http"
	"fmt"
)

/*
HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 13

Hello World!
*/
type Server struct {
	listener net.Listener
	isClosed atomic.Bool
}

// stops the server by closing the underlying net.Listener. 
// Setting the atomic boolean ensures the listen() loop 
// knows the shutdown was intentional
func (s *Server) Close() error {
	s.isClosed.Store(true) // Mark as closed before closing the listener
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// runs the acceptance loop. By checking the atomic.Bool, 
// you can distinguish between a real network error and 
// an expected error caused by calling Close()
func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
				// If the server was intentionally closed, ignore the error and exit
			if s.isClosed.Load() {
				return
			}
			fmt.Printf("Accept error: %v\n", err)
			continue
		}

		// Handle each connection in its own goroutine
		go s.handle(conn)
	}
}

// manages the lifecycle of a single connection. 
// It is critical to use defer conn.Close() to ensure 
// the TCP connection is released regardless of how the function exits. 
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	//example response
	response := "HTTP/1.1 200 OK\r\n" +
							"Content-Type: text/plain\r\n" +
							"Content-Length: 13\r\n" +
							"\r\n" +
							"Hello, World!"
	
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
	}
}

func Serve(port int) (*Server, error) {
    addr := fmt.Sprintf(":%d", port)
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        return nil, err
    }

    srv := &Server{
        listener: ln,
        // other configuration
    }

    go srv.listen()

    return srv, nil
}
