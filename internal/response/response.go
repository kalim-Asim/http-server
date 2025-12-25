package response

import (
	"fmt"
	"io"
	"github.com/kalim-Asim/http-server/internal/headers"
)

type StatusCode int 

type Response struct {
}

const (
	StatusOK StatusCode = 200 
	StatusBadRequest StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var line string

	switch statusCode {
	case StatusOK: // 200
		line = "HTTP/1.1 200 OK\r\n"
	case StatusBadRequest: // 400
		line = "HTTP/1.1 400 Bad Request\r\n"
	case StatusInternalServerError: // 500
		line = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		// Any other code leaves the reason phrase blank
		line = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
	}

	_, err := w.Write([]byte(line))
	return err
}

// it should set the following headers that we always want to include in our responses
func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("content-length", fmt.Sprintf("%d", contentLen))
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")
	return *h 
}

// add our responses to header
func WriteHeaders(w io.Writer, h headers.Headers) error {
	b := []byte{} 

	h.ForEach(func(key, val string){
		b = fmt.Appendf(b, "%s: %s\r\n", key, val)
	})

	b = fmt.Appendf(b, "\r\n")
	_, err := w.Write(b)

	return err  
}