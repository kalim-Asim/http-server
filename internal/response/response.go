package response

import (
	"fmt"
	"io"
	"github.com/kalim-Asim/http-server/internal/headers"
)

type StatusCode int 
const (
	StatusOK StatusCode = 200 
	StatusBadRequest StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type Response struct {
}

// it should set the following headers that we always want to include in our responses
func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h 
}

type Writer struct {
	writer io.Writer 
}

func NewWriter(w io.Writer) *Writer{
	return &Writer{
		writer: w, 
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	var line string

	switch statusCode {
		case StatusOK: 
			line = "HTTP/1.1 200 OK\r\n"
		case StatusBadRequest: 
			line = "HTTP/1.1 400 Bad Request\r\n"
		case StatusInternalServerError: 
			line = "HTTP/1.1 500 Internal Server Error\r\n"
		default:
			// Any other code leaves the reason phrase blank
			line = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
	}

	_, err := w.writer.Write([]byte(line))
	return err
}
 
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	b := []byte{} 

	headers.ForEach(func(key, val string){
		b = fmt.Appendf(b, "%s: %s\r\n", key, val)
	})

	b = fmt.Appendf(b, "\r\n")
	_, err := w.writer.Write(b)

	return err  
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	return n, err 
}

/*
HTTP/1.1 200 OK
Content-Type: text/plain
Transfer-Encoding: chunked

<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
<n>\r\n
<data of length n>\r\n
... repeat ...
0\r\n
\r\n



Example:- 
HTTP/1.1 200 OK
Content-Type: text/plain
Transfer-Encoding: chunked
Trailer: Lane, Prime, TJ

1E
I could go for a cup of coffee
C
But not Java
12
Never go full Java
0

0\r\n
Lane: goober
Prime: chill-guy
TJ: 1-indexer
\r\n
*/

// func (w *Writer) WriteChunkedBody(p []byte) (int, error) {

// }

// func (w *Writer) WriteChunkedBodyDone() (int, error) {

// }

// func (w *Writer) WriteTrailers(t *headers.Headers, body []byte) error {

// }

