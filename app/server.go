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
type Middleware func(RequestHandler) RequestHandler

type Server struct {
	Addr        int
	Handler     map[string]map[string]RequestHandler
	Middlewares []Middleware
}

func NewServer(port int) *Server {
	return &Server{
		Addr: port,
		Handler: map[string]map[string]RequestHandler{
			"GET":    {},
			"POST":   {},
			"PUT":    {},
			"DELETE": {},
		},
	}
}

func (s *Server) Get(path string, handler RequestHandler) {
	s.Handler["GET"][path] = handler
}

func (s *Server) Post(path string, handler RequestHandler) {
	s.Handler["POST"][path] = handler
}

func (s *Server) Use(middleware Middleware) {
	s.Middlewares = append(s.Middlewares, middleware)
}

func chainMiddlewares(h RequestHandler, middlewares ...Middleware) RequestHandler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.Addr))
	if err != nil {
		fmt.Printf("Failed to bind to port %d", s.Addr)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)

				if err != nil {
					if err != io.EOF {
						fmt.Println("Read error:", err)
					}
					return
				}

				data := buf[:n]
				request, err := NewRequest(data)
				if err != nil {
					conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
				}

				var response *Response
				for route, handler := range s.Handler[request.Method] {
					re := regexp.MustCompile(route)

					if re.Match([]byte(request.Path)) {
						finalHandler := chainMiddlewares(handler, s.Middlewares...)
						response = finalHandler(request)
						conn.Write([]byte(response.Content(request.HTTPVersion)))
						return
					}
				}

				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

				if request.GetHeader("Connection") == "close" {
					return
				}
			}
		}(conn)
	}
}
