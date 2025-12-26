package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kalim-Asim/http-server/internal/request"
	"github.com/kalim-Asim/http-server/internal/response"
	"github.com/kalim-Asim/http-server/internal/server"
)

const port = 42069
const BadRequest = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
const InternalServerError = `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`
const StatusOk = `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

func main() {
	server, err := server.Serve(
		port,
		func(w *response.Writer, req *request.Request) {
			h := response.GetDefaultHeaders(0)
			body := []byte("All good, frfr\n")
			var status response.StatusCode = response.StatusOK

			if req.RequestLine.RequestTarget == "/yourproblem" {
				body = []byte(BadRequest)
				status = response.StatusBadRequest
			} else if req.RequestLine.RequestTarget == "/myproblem" {
				body = []byte(InternalServerError)
				status = response.StatusInternalServerError
			} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream") {
				target := req.RequestLine.RequestTarget
				res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
				if err != nil {
					body = []byte(InternalServerError)
					status = response.StatusInternalServerError
				} else {
					w.WriteStatusLine(status)
					slog.Info("https://httpbin.org/" + target[len("/httpbin/"):])
					h.Delete("Content-Length")
					h.Set("Transfer-Encoding", "chunked")
					h.Set("Content-Type", "text/plain")
					w.WriteHeaders(*h)

					for {
						data := make([]byte, 16)
						n, err := res.Body.Read(data)
						if err != nil {
							break
						}

						w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
						w.WriteBody(data[:n])
						w.WriteBody([]byte("\r\n"))
					}
					w.WriteBody([]byte("0\r\n\r\n"))
				}

			}

			h.Set("Content-Length", fmt.Sprintf("%d", len(body)))
			h.Set("Content-Type", "text/html")
			w.WriteStatusLine(status)
			w.WriteHeaders(*h)
			w.WriteBody(body)
		})

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
