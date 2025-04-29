package main

import (
	"bytes"
	"errors"
	"strings"
)

type Request struct {
	Method      string
	Path        string
	HTTPVersion string

	Headers map[string][]string
}

func NewRequest(content []byte) (*Request, error) {
	requestParts := bytes.Split(content, []byte("\r\n"))

	requestLine := requestParts[0]
	requestHeaders := requestParts[1 : len(requestParts)-1]

	requestLineParts := bytes.Split(requestLine, []byte(" "))
	if len(requestLineParts) != 3 {
		return nil, errors.New("Could not parse request line")
	}

	method := requestLineParts[0]
	url := requestLineParts[1]
	version := requestLineParts[2]
	headers := parseHeaders(requestHeaders)

	return &Request{
		Method:      string(method),
		Path:        string(url),
		HTTPVersion: string(version),
		Headers:     headers,
	}, nil
}

func parseHeaders(headersLines [][]byte) map[string][]string {
	headers := make(map[string][]string)

	for _, headerLine := range headersLines {
		splittedHeaderLine := bytes.SplitN(headerLine, []byte(":"), 2)
		if len(splittedHeaderLine) != 2 {
			continue
		}

		header := string(bytes.ToLower(bytes.TrimSpace(splittedHeaderLine[0])))
		headerValue := string(bytes.TrimSpace(splittedHeaderLine[1]))

		headers[header] = append(headers[header], headerValue)
	}

	return headers
}

func (r *Request) GetHeader(header string) string {
	normalizedHeader := strings.ToLower(strings.TrimSpace(header))

	value := r.Headers[normalizedHeader]
	return strings.Join(value, ", ")
}
