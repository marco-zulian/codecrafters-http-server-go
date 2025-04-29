package main

import (
	"slices"
	"strings"
)

func EncodingMiddleware(next RequestHandler) RequestHandler {
	return func(r *Request) *Response {
		response := next(r)
		acceptedEncodings := strings.Split(r.GetHeader("Accept-Encoding"), ", ")

		if slices.Contains(acceptedEncodings, "gzip") {
			response.AddHeader("Content-Encoding", "gzip")
		}

		return response
	}
}
