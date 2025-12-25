package request

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/kalim-Asim/http-server/internal/headers"
)

// we have to parse this...
/*
	POST /coffee HTTP/1.1
	Host: localhost:42069
	User-Agent: curl/7.81.0
	Accept: * / *
	Content-Length: 21

	{"flavor":"dark mode"}
*/

/*
HTTP-message   = start-line CRLF
							*( field-line CRLF )
							CRLF
							[ message-body ] ( we need to parse this now..)
*/

var (
	SEPARATOR = []byte("\r\n")
	ERROR_BAD_START_LINE = fmt.Errorf("bad start line")
	ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
	ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("http version not supported")
)

// parser state machine, to track parser progress
type parserState string 
const (
	StateInit parserState = "init"
	StateHeader parserState = "header"
	StateBody parserState = "body"
	StateDone parserState = "done"
	StateError parserState = "error"
)

// example: request-line -> GET /index.html HTTP/1.1
// Represents the first line of an HTTP request
type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

// Represents the parsed request so far
type Request struct {
	RequestLine RequestLine // holds the parse requestline(first line)
	State parserState
	Headers headers.Headers // headers parsed
	Body string
}

func NewRequest() *Request {
	return &Request{
		State: StateInit,
		Headers: *headers.NewHeaders(), 
	}
}

func (r *Request) done() bool {
	return r.State == StateDone
}

// helper function to do parsing
func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil 
	}

	startLine, read := b[:idx], idx + len(SEPARATOR)
	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_BAD_START_LINE
	}
	
	httpParts := bytes.Split(parts[2], []byte("/"))
	// valid http
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_BAD_START_LINE
	}

	rl := &RequestLine{
		Method: string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion: string(httpParts[1]),
	} 

	return rl, read, nil 
}

func getLength(r headers.Headers, key string) int {
	var length int = 0 
	if r.Has(key) {
		length, _ = strconv.Atoi(r.Get(key))
	}
	return length 
}
func (r *Request) hasBody() bool {
	length := getLength(r.Headers, "content-length")
	return length > 0 
}

// It accepts the next slice of bytes that needs to be parsed into the Request struct.
// It updates the "state" of the parser, and the parsed RequestLine field.
// It returns the number of bytes it consumed 
// (meaning successfully parsed)
func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break 
		}

		switch r.State {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE 
			
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				return 0, err 
			}
			if n == 0 {
				break outer 
			}
			
			r.RequestLine = *rl 
			read += n 
			
			r.State = StateHeader 
			
		case StateHeader:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.State = StateError
				return 0, fmt.Errorf("error parsing header... ")
			}
			if n == 0 {
				break outer
			}
			read += n 

			if done {
				if r.hasBody() {
					r.State = StateBody
				} else {
					r.State = StateDone
				}
			}

		case StateBody:

			length := getLength(r.Headers, "content-length")
			if length == 0 {
					panic("chuncked not implemented")
			}

			remaining := min(len(currentData), length - len(r.Body))
			r.Body += string(currentData[:remaining])
			read += remaining

			if len(r.Body) == length {
					r.State = StateDone
			}

		case StateDone:
			break outer 

		default:
			panic("nothing to show... ")
		}
	}
	
	return read, nil 
}

// orchestration function,
// parse the request-line from the reader
func RequestFromReader(reader io.Reader) (*Request, error) {
	req := NewRequest()
	buf := make([]byte, 1024)
	bufLen := 0

	for !req.done() {
		slog.Info("RequestFromReader", "state", req.State)

		n, err := reader.Read(buf[bufLen:])
		if err != nil {
				return nil, err
		}

		bufLen += n

		readN, parseErr := req.parse(buf[:bufLen])
		if parseErr != nil {
				return nil, parseErr
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return req, nil 	
}