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
	server.Get("^/$", EmptyPathHandler)
	server.Get("^/echo", EchoPathHandler)
	server.Get("^/user-agent", UserAgentHandler)
	server.Get("^/files", FilesGetPathHandler)
	server.Post("^/files/.+", FilesPostPathHandler)

	server.Use(EncodingMiddleware)

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

func FilesGetPathHandler(request *Request) *Response {
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

func FilesPostPathHandler(request *Request) *Response {
	response := NewResponse()

	filePathParts := strings.SplitN(request.Path, "/", 3)

	f, err := os.Create(fmt.Sprintf("%s%s", *directory, filePathParts[2]))
	if err != nil {
		response.SetStatus(500)
		return response
	}
	defer f.Close()

	_, err = f.WriteString(request.Body)
	if err != nil {
		response.SetStatus(500)
		return response
	}

	response.SetStatus(201)
	return response
}
