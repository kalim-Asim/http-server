package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var SEPARATOR = []byte("\r\n") // cflf token
var ERROR_CRLF_NOT_FOUND = fmt.Errorf("crlf token not found")
var ERROR_BAD_HEADER = fmt.Errorf("header does not match")

// example: header -> Host: localhost:42069\r\n\r\n
type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers) 
}

// there should not be any space between key and first colon
func isValidHeader(d []byte) (bool, string, string) {
	s := strings.TrimSpace(string(d))

	idxColon := strings.IndexByte(s, ':')
	if idxColon == -1 {
		return false, "", ""
	}

	// No space before colon
	if idxColon > 0 && s[idxColon-1] == ' ' {
		return false, "", ""
	}

	key := strings.TrimSpace(s[:idxColon])
	value := strings.TrimSpace(s[idxColon+1:])

	if key == "" || value == "" {
		return false, "", ""
	}

	return true, key, value
}


// parse header line by line and adds into our map
// done=true when crlf is at start
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, SEPARATOR)
	if idx == -1 {
		return 0, false, nil
	}

	// end of headers
	if idx == 0 {
		return len(SEPARATOR), true, nil
	}

	line, read := data[:idx], idx + len(SEPARATOR)

	isValid, key, value := isValidHeader(line)
	if !isValid {
		return 0, false, ERROR_BAD_HEADER
	} 

	h[key] = value 
	return read, false, nil 
} 

