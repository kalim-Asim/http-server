package request

import (
	"bytes"
	"fmt"
	"io"
)

// focused only on parsing the HTTP request line

// parser state machine, to track parser progress
type parserState string 
const (
	StateInit parserState = "init"
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
	RequestLine RequestLine // holds the parse requestline
	State parserState
}

func NewRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad start line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("http version not supported")
var ERROR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n") 

func (r *Request) done() bool {
	return r.State == StateDone
}

func (r *Request) error() bool {
	return r.State == StateError
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

// It accepts the next slice of bytes that needs to be parsed into the Request struct.
// It updates the "state" of the parser, and the parsed RequestLine field.
// It returns the number of bytes it consumed 
// (meaning successfully parsed)
func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.State {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATE 

		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err 
			}
			if n == 0 {
				break outer 
			}

			r.RequestLine = *rl 
			read += n 

			r.State = StateDone 

		case StateDone:
			break outer 
		}
	}
	
	return read, nil 
}

// orchestration function
// parse the request-line from the reader
func RequestFromReader(reader io.Reader) (*Request, error) {
	req := NewRequest()
	buf := make([]byte, 1024)
	bufLen := 0

	for !req.done() && !req.error(){
		n, err := reader.Read(buf[bufLen:]) // how many bytes you have read from the reader
		if err != nil {
			return nil, err 
		}

		bufLen += n
		readN, err := req.parse(buf[:bufLen]) // how many bytes you have parsed from the buffer
		if err != nil {
			return nil, err 
		}
		
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return req, nil 	
}