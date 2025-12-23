package headers

import (
	"bytes"
	"fmt"
)

var SEPARATOR = []byte("\r\n") // cflf token
var ERROR_CRLF_NOT_FOUND = fmt.Errorf("crlf token not found")
var ERROR_BAD_HEADER = fmt.Errorf("header does not match")
var ERROR_BAD_FIELD_NAME = fmt.Errorf("malformed field name")

// example: header -> Host: localhost:42069\r\n\r\n
type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers) 
}

// returns key, value, error 
func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ERROR_BAD_FIELD_NAME
	}

	rawName := parts[0]
	if len(rawName) > 0 && (rawName[len(rawName)-1] == ' ' || rawName[len(rawName)-1] == '\t') {
		return "", "", ERROR_BAD_FIELD_NAME
	}

	name := bytes.TrimSpace(rawName)
	val := bytes.TrimSpace(parts[1])

	return string(name), string(val), nil
}

// parse header line by line and adds into our map
// done=true when crlf is at start
func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], SEPARATOR)
		if idx == -1 {
			break // no more complete lines
		}

		// check for empty line -> headers done
		if idx == 0 {
			if read == 0 {
        // input starts with CRLF
        read += len(SEPARATOR)
				done = true
    	}
			break
		}

		line := data[read : read+idx]
		fmt.Printf("%s\n", string(line))

		key, val, err := parseHeader(line)
		if err != nil {
			return 0, false, err
		}

		h[key] = val
		read += idx + len(SEPARATOR)
	}

	return read, done, nil
}