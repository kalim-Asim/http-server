package headers

import (
	"bytes"
	"fmt"
	"strings"
)

/* -------------  HEADER FORMAT  ----------------

field-line   = field-name ":" OWS field-value OWS
*/

var (
	SEPARATOR = []byte("\r\n") // cflf token
	ERROR_BAD_FIELD_NAME = fmt.Errorf("malformed field name")
	ERROR_CRLF_NOT_FOUND = fmt.Errorf("crlf token not found")
	ERROR_BAD_HEADER = fmt.Errorf("header does not match")
	ERROR_INVALID_FIELD_NAME = fmt.Errorf("field name is invalid")
)

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

// example: header -> Host: localhost:42069\r\n (valid)
type Headers struct {
	// key lower case
	// value can be multiple as per RFC 9110 5.2
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: make(map[string]string),
	}
}

func (h *Headers) Get(key string) string {
	return h.headers[strings.ToLower(key)]
}

func (h* Headers) Has(key string) bool {
	_, ok := h.headers[strings.ToLower(key)]
	return ok 
}

func (h *Headers) Set(key, value string) {
	h.headers[strings.ToLower(key)] = value 
}

func (h* Headers) PrintAll() {
	for key, val := range h.headers {
		fmt.Printf(" - %s: %s\n", key, val)
	}
}
// parse header line by line and adds into our map
// done=true when crlf is at start
func (h* Headers) Parse(data []byte) (int, bool, error) {
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
	
		key, val, err := parseHeader(line)
		if err != nil {
			return 0, false, err
		}
		if !isToken(key) {
			return 0, false, ERROR_INVALID_FIELD_NAME
		}

		if h.Has(key) {
			oldValue := h.Get(key)
			newValue := oldValue + ", " + val 
			h.Set(key, newValue)
		} else {
			h.Set(key, val)
		}
		
		read += idx + len(SEPARATOR)
	}

	return read, done, nil
}

func isToken(str string) bool {
	for _, ch := range str  {
		found := false 
		if ch >= 'A' && ch <= 'Z'|| 
			ch >= 'a' && ch <= 'z'|| 
			ch >= '0' && ch <= '9' {
				found = true 
			}

		switch ch {
			case  '#', '!', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '/', '~':
				found = true 
		}

		if !found {
			return false
		}
	}
	return len(str) >= 1 
}