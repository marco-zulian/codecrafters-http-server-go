package main

import (
	"errors"
	"fmt"
	"strings"
)

var httpStatusCodes = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Content Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	421: "Misdirect Request",
	422: "Unprocessable Content",
	423: "Locked",
	424: "Failed Dependency",
	425: "Too Early",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}

type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       string
}

func NewResponse() *Response {
	return &Response{
		StatusCode: 200,
		Headers:    make(map[string][]string),
		Body:       "",
	}
}

func (r *Response) SetStatus(status int) error {
	_, validStatus := httpStatusCodes[status]

	if !validStatus {
		return errors.New("Unknown status code")
	}

	r.StatusCode = status
	return nil
}

func (r *Response) AddHeader(header string, value string) {
	normalizedHeader := strings.ToLower(strings.TrimSpace(header))

	r.Headers[normalizedHeader] = append(r.Headers[normalizedHeader], value)
}

func (r *Response) SetHeader(header string, value string) {
	normalizedHeader := strings.ToLower(strings.TrimSpace(header))

	r.Headers[normalizedHeader] = []string{value}
}

func (r *Response) SetBody(body string) {
	r.Body = body
}

func (r *Response) GetHeaders() string {
	var b strings.Builder

	for header, value := range r.Headers {
		b.WriteString(fmt.Sprintf("%s: %s\r\n", header, strings.Join(value, ", ")))
	}

	return b.String()
}

func (r *Response) Content(protocol string) string {
	return fmt.Sprintf("%s %d %s\r\n%s\r\n%s", protocol, r.StatusCode, httpStatusCodes[r.StatusCode], r.GetHeaders(), r.Body)
}
