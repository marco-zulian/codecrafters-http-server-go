package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	server := NewServer(4221)
	server.AddHandler("^/$", EmptyPathHandler)
	server.AddHandler("^/echo", EchoPathHandler)
  server.AddHandler("^/user-agent", UserAgentHandler)

	server.Serve()
}

func EmptyPathHandler(request *Request) *Response {
	return NewResponse()
}

func EchoPathHandler(request *Request) *Response {
	response := NewResponse()

	echoPathParts := strings.SplitN(request.Path, "/", 3)
  
  var echoStr string
  if len(echoPathParts) > 2 {
    echoStr = echoPathParts[2]
  }

	response.AddHeader("Content-Type", "text/plain")
	response.AddHeader("Content-Length", strconv.Itoa(len(echoStr)))
	response.SetBody(echoStr)

	return response
}

func UserAgentHandler(request *Request) *Response {
  response := NewResponse()

  userAgent := request.GetHeader("User-Agent")
  
  response.AddHeader("Content-Type", "text/plain")
  response.AddHeader("Content-Length", strconv.Itoa(len(userAgent)))
  response.SetBody(userAgent)

  return response
}
