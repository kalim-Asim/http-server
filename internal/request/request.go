package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	Method        string
	RequestTarget string
	HttpVersion   string
}

// full parsed HTTP request
type Request struct {
	RequestLine RequestLine
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad start line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("http version not supported")
var SEPARATOR = "\r\n"

// helper function to do parsing
func parseRequestLine(s string) (*RequestLine, string, error) {
	idx := strings.Index(s, SEPARATOR)
	if idx == -1 {
		return nil, s, nil 
	}

	startLine, restOfmsg := s[:idx], s[idx + len(SEPARATOR):]
	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfmsg, ERROR_BAD_START_LINE
	}
	

	httpParts := strings.Split(parts[2], "/")
	// valid http
	if len(httpParts) != 2 || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil, restOfmsg, ERROR_BAD_START_LINE
	}

	rl := &RequestLine{
		Method: parts[0],
		RequestTarget: parts[1],
		HttpVersion: httpParts[1],
	} 

	return rl, restOfmsg, nil 
}

//  parse the request-line from the reader
func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("unable to io.ReadAll"),
			err, 
		)
	}

	str := string(data)
	rl, _, err := parseRequestLine(str)
	if err != nil {
		return nil, err 
	}

	return &Request{
		RequestLine: *rl,
	}, err 
}