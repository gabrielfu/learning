package main

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func (h Headers) String() string {
	var result string
	for name, value := range h {
		result += fmt.Sprintf("%s: %s\r\n", name, value)
	}
	return result
}

type Request struct {
	Method      string
	Path        string
	HttpVersion string
	Headers     Headers
	Body        []byte
}

func parseRequest(value string) (*Request, error) {
	segments := strings.Split(value, "\r\n")
	if len(segments) < 1 {
		return nil, fmt.Errorf("invalid request: %v", value)
	}
	requestLine := strings.Split(segments[0], " ")
	if len(requestLine) < 3 {
		return nil, fmt.Errorf("invalid request line: %v", requestLine)
	}
	method := requestLine[0]
	path := requestLine[1]
	httpVersion := requestLine[2]

	headers := make(Headers)
	for _, header := range segments[1 : len(segments)-1] {
		if header == "" {
			break
		}
		headerParts := strings.SplitN(header, ":", 2)
		if len(headerParts) != 2 {
			return nil, fmt.Errorf("invalid header: %v", header)
		}
		name := strings.TrimSpace(headerParts[0])
		value := strings.TrimSpace(headerParts[1])
		headers[name] = value
	}
	body := []byte(segments[len(segments)-1])
	return &Request{
		Method:      method,
		Path:        path,
		HttpVersion: httpVersion,
		Headers:     headers,
		Body:        body,
	}, nil
}

type Response struct {
	HttpVersion   string
	StatusCode    int
	StatusMessage string
	Headers       Headers
	Body          []byte
}

func (r Response) String() string {
	return fmt.Sprintf("%s %d %s\r\n%s\r\n%s", r.HttpVersion, r.StatusCode, r.StatusMessage, r.Headers.String(), string(r.Body))
}
