package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var directory = flag.String("directory", "", "Specify the directory from where to search for files")

func main() {
	fmt.Println("Logs from your program will appear here!")
	flag.Parse()

	server := NewServer(4221)
	server.AddHandler("^/$", EmptyPathHandler)
	server.AddHandler("^/echo", EchoPathHandler)
	server.AddHandler("^/user-agent", UserAgentHandler)
	server.AddHandler("^/files", FilesPathHandler)

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

func FilesPathHandler(request *Request) *Response {
	response := NewResponse()

	filePathParts := strings.SplitN(request.Path, "/", 3)

	if len(filePathParts) != 3 {
		response.SetStatus(404)
		return response
	}

	data, err := os.ReadFile(fmt.Sprintf("%s%s", *directory, filePathParts[2]))
	if err != nil {
		fmt.Println(err)
		response.SetStatus(404)
		return response
	}

	response.SetBody(string(data))
	response.AddHeader("Content-Type", "application/octet-stream")
	response.AddHeader("Content-Length", strconv.Itoa(len(data)))

	return response
}
