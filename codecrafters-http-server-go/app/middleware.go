package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"slices"
	"strings"
)

type Middleware interface {
	Dispatch(request Request, callNext func(Request) Response) Response
}

type GzipEncodingMiddleware struct{}

func (m GzipEncodingMiddleware) Dispatch(request Request, callNext func(Request) Response) Response {
	acceptEncodings := strings.Split(request.Headers["Accept-Encoding"], ",")
	for i := range acceptEncodings {
		acceptEncodings[i] = strings.TrimSpace(acceptEncodings[i])
	}

	response := callNext(request)
	if !slices.Contains(acceptEncodings, "gzip") {
		return response
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(response.Body)
	if err != nil {
		return Response{
			HttpVersion:   request.HttpVersion,
			StatusCode:    500,
			StatusMessage: "Internal Server Error",
			Headers:       Headers{},
			Body:          []byte("Internal Server Error"),
		}
	}
	gz.Close()
	response.Body = b.Bytes()
	response.Headers["Content-Encoding"] = "gzip"
	response.Headers["Content-Length"] = fmt.Sprintf("%d", len(response.Body))
	return response
}
