package main

import (
	"fmt"
	"net"
  "os"
  "io"
  "bytes"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
	 	fmt.Println("Failed to bind to port 4221")
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
    requestLine := bytes.Split(data, []byte("\r\n"))[0]
    path := bytes.TrimPrefix(bytes.Split(requestLine, []byte(" "))[1], []byte("/"))
    
    
    fmt.Println(string(path))
    if (string(path) == "") {
      conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
      return
    }

    if (bytes.HasPrefix(path, []byte("echo"))) {
      echoStr := bytes.Split(path, []byte("/"))[1]
      conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoStr), echoStr)))
      return
    }
    

    conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
    break
  }
}

