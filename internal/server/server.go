package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"github.com/kalim-Asim/http-server/internal/request"
	"github.com/kalim-Asim/http-server/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
	handler  Handler
}

//a proper status code and error message
type HandlerError struct {
	StatusCode   response.StatusCode 
	Message []byte 
}

type Handler func(w io.Writer, req *request.Request) *HandlerError 

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

// manages the lifecycle of a single connection. 
// It is critical to use defer conn.Close() to ensure 
// the TCP connection is released regardless of how the function exits. 
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)

	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, *headers)
		return 
	}

	var body []byte = nil 
	writer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(writer, r)
	var status response.StatusCode = response.StatusOK

	if handlerError != nil {
		status = handlerError.StatusCode
		body = []byte(handlerError.Message)
	} else {
		body = writer.Bytes()
	}

	headers.Set("Content-Length", fmt.Sprintf("%d", len(body)))
	response.WriteStatusLine(conn, status)
	response.WriteHeaders(conn, *headers)
	conn.Write(body)
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

		go s.handle(conn)
	}
}

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
			return nil, err
	}

	srv := &Server{
		listener: ln,
		handler: handler,
	}

	go srv.listen()

	return srv, nil
}