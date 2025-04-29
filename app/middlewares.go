package main

import (
	"bytes"
	"compress/gzip"
	"slices"
	"strconv"
	"strings"
)

func EncodingMiddleware(next RequestHandler) RequestHandler {
	return func(r *Request) *Response {
		response := next(r)
		acceptedEncodings := strings.Split(r.GetHeader("Accept-Encoding"), ", ")

		if slices.Contains(acceptedEncodings, "gzip") {
			var buf bytes.Buffer
			gw := gzip.NewWriter(&buf)

			_, err := gw.Write([]byte(r.Body))
			if err != nil {
				return response
			}

			err = gw.Close()
			if err != nil {
				return response
			}

			compressedBody := buf.String()

			response.SetBody(compressedBody)
			response.SetHeader("Content-Length", strconv.Itoa(len(compressedBody)))
			response.AddHeader("Content-Encoding", "gzip")
		}

		return response
	}
}
