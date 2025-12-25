package main

import (
	// "io"
	"fmt"
	"log"
	"os"
	"os/signal"
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

			switch req.RequestLine.RequestTarget {
				case "/yourproblem":
					body = []byte(BadRequest)
					status = response.StatusBadRequest

				case "/myproblem":
					body = []byte(InternalServerError)
					status = response.StatusInternalServerError 
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