package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
)

var _ = net.Listen
var _ = os.Exit

type RequestHandler func(*Request) *Response

type Server struct {
	Addr    int
	Handler map[string]RequestHandler
}

func NewServer(port int) *Server {
	return &Server{
		Addr:    port,
		Handler: make(map[string]RequestHandler),
	}
}

func (s *Server) AddHandler(path string, handler RequestHandler) {
	s.Handler[path] = handler
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.Addr))
	if err != nil {
		fmt.Printf("Failed to bind to port %d", s.Addr)
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			break
		}

		data := buf[:n]
		request, err := NewRequest(data)
		if err != nil {
			return err
		}

		var response *Response
		for route, handler := range s.Handler {
			re := regexp.MustCompile(route)

			if re.Match([]byte(request.Path)) {
				response = handler(request)
				conn.Write([]byte(response.Content(request.HTTPVersion)))
				return nil
			}
		}

		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		break
	}

	return nil
}
