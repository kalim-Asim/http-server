# HTTP from TCP

A minimal HTTP/1.1 server implemented directly on top of TCP using Go.  
This project avoids `net/http` for serving and instead focuses on understanding HTTP at the wire level: request parsing, headers, bodies, chunked transfer encoding, trailers, and correct message framing per RFC 9112.

---

## Features

- Raw TCP-based HTTP/1.1 server
- Manual parsing of:
  - Request line
  - Headers (RFC-compliant field names and values)
  - Message body
- Proper response formatting (status line, headers, body)
- Chunked Transfer-Encoding with trailers
- Streaming responses
- Binary-safe responses (video)
- Debug TCP listener for inspecting raw requests

---

## Endpoints

| Method | Path | Description |
|------|------|------------|
| GET | `/` | `200 OK` — `All good, frfr\n` |
| GET | `/yourproblem` | `400 Bad Request` — `Your problem is not my problem\n` |
| GET | `/myproblem` | `500 Internal Server Error` — `Woopsie, my bad\n` |
| GET | `/video` | Serves `assets/vim.mp4` with `video/mp4` |
| GET | `/httpbin/stream/100` | Raw chunked response streamed byte-by-byte |

---

## Project Structure

```text
http-server/
├── assets/
│   └── vim.mp4              # Video file served by /video endpoint(make sure to add it)
│
├── cmd/
│   ├── httpserver/
│   │   └── main.go          # Main TCP HTTP server entry point
│   │
│   ├── tcplistener/
│   │   └── main.go          # Raw TCP listener for debugging requests
│   │
│   └── udpsender/
│       └── main.go          # Simple UDP sender
│
├── internal/
│   ├── headers/
│   │   ├── headers.go       # HTTP header storage and parsing logic
│   │   └── headers_test.go  # Unit tests for headers
│   │
│   ├── request/
│   │   ├── request.go       # HTTP request parsing from TCP stream
│   │   └── request_test.go  # Request parsing tests
│   │
│   ├── response/
│   │   └── response.go     # HTTP response writer (status, headers, body, chunked)
│   │
│   └── server/
│       └── server.go        # TCP server accept loop and routing
│
├── messages.txt             # Test / sample HTTP messages(did in starting)
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── .gitignore
└── Readme.md
````

### Design Notes

* `cmd/` contains executable entry points only
* `internal/` holds all protocol logic (headers, request parsing, response writing)
* `net/http` is used only for proxing to httpbin.org
* HTTP/1.1 framing, chunked encoding, and trailers are handled manually
---

## Chunked Transfer Encoding

The `/httpbin/stream/100` endpoint proxies `httpbin.org` and returns a **raw chunked response**, including:

- Hexadecimal chunk sizes
- CRLF delimiters
- Final zero-length chunk
- HTTP trailers:
  - `X-Content-SHA256`
  - `X-Content-Length`

### Important Note

Use `curl --raw` or `nc` to prevent the client from decoding the chunked response.

---

## Verifying Raw Chunked Output

```bash
echo -e "GET /httpbin/stream/100 HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n" | nc localhost 42069
````

This shows the exact bytes sent over the wire, including chunk sizes and trailers.

---

## Running the Server

From the project root:

```bash
go run ./cmd/httpserver
```

Make sure the video file exists:

```text
assets/vim.mp4
```

---

## TCP Listener (Debug Tool)

A separate TCP listener is provided to inspect incoming requests exactly as received.

Run the listener and redirect output to a file:

```bash
go run ./cmd/tcplistener | tee /tmp/body.txt
```

In another shell, send a request:

```bash
curl -X POST http://localhost:42069/coffee \
  -H 'Content-Type: application/json' \
  -d '{"type":"dark mode","size":"medium"}'
```

This is useful for validating request parsing and body handling.

---

## HTTP Parsing Notes

### Message Format

```
HTTP-message = start-line CRLF
               *( field-line CRLF )
               CRLF
               [ message-body ]
```

### Request Line

* Method
* Request Target
* HTTP Version

### Headers

```
field-line = field-name ":" OWS field-value OWS
```

* Header names are case-insensitive
* No whitespace allowed in field names
* Unlimited optional whitespace around field values
* Field names must be valid RFC tokens

### Responses

Responses follow the same structure as requests, except the start-line is a status line:

```
status-line = HTTP-version SP status-code SP [ reason-phrase ]
```

Common response headers used:

* `Content-Length`
* `Content-Type`
* `Connection`
* `Transfer-Encoding`
* `Trailer`

---

## Motivation

This project exists to gain a deep, practical understanding of HTTP by implementing it manually over TCP, especially areas often abstracted away such as chunked encoding, trailers, and strict protocol correctness.

---

