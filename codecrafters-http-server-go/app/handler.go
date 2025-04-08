package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

var directoryFlag string

type Handler func(Request) *Response

var handlers = []Handler{
	home,
	echo,
	userAgent,
	getFile,
	postFile,
}

func handleRequest(r Request) Response {
	for _, handler := range handlers {
		if resp := handler(r); resp != nil {
			return *resp
		}
	}

	return Response{
		HttpVersion:   r.HttpVersion,
		StatusCode:    404,
		StatusMessage: "Not Found",
		Headers:       Headers{},
		Body:          []byte{},
	}
}

func home(r Request) *Response {
	if r.Path == "/" {
		return &Response{
			HttpVersion:   r.HttpVersion,
			StatusCode:    200,
			StatusMessage: "OK",
			Headers:       Headers{},
			Body:          []byte{},
		}
	}
	return nil
}

func echo(r Request) *Response {
	if strings.HasPrefix(r.Path, "/echo/") {
		echoStr := strings.TrimPrefix(r.Path, "/echo/")
		return &Response{
			HttpVersion:   r.HttpVersion,
			StatusCode:    200,
			StatusMessage: "OK",
			Headers: Headers{
				"Content-Type":   "text/plain",
				"Content-Length": fmt.Sprintf("%d", len(echoStr)),
			},
			Body: []byte(echoStr),
		}
	}
	return nil
}

func userAgent(r Request) *Response {
	if r.Path == "/user-agent" {
		userAgent := r.Headers["User-Agent"]
		return &Response{
			HttpVersion:   r.HttpVersion,
			StatusCode:    200,
			StatusMessage: "OK",
			Headers: Headers{
				"Content-Type":   "text/plain",
				"Content-Length": fmt.Sprintf("%d", len(userAgent)),
			},
			Body: []byte(userAgent),
		}
	}
	return nil
}

func getFile(r Request) *Response {
	if strings.HasPrefix(r.Path, "/files/") && r.Method == "GET" {
		filePath := strings.TrimPrefix(r.Path, "/files/")
		if strings.Contains(filePath, "..") {
			return &Response{
				HttpVersion:   r.HttpVersion,
				StatusCode:    400,
				StatusMessage: "Bad Request",
				Headers:       Headers{},
				Body:          []byte("Invalid file path"),
			}
		}

		fullPath := path.Join(directoryFlag, filePath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return &Response{
					HttpVersion:   r.HttpVersion,
					StatusCode:    404,
					StatusMessage: "Not Found",
					Headers:       Headers{},
					Body:          []byte("File not found"),
				}
			}
			return &Response{
				HttpVersion:   r.HttpVersion,
				StatusCode:    500,
				StatusMessage: "Internal Server Error",
				Headers:       Headers{},
				Body:          []byte("Error reading file"),
			}
		}
		return &Response{
			HttpVersion:   r.HttpVersion,
			StatusCode:    200,
			StatusMessage: "OK",
			Headers: Headers{
				"Content-Type":   "application/octet-stream",
				"Content-Length": fmt.Sprintf("%d", len(content)),
			},
			Body: content,
		}
	}
	return nil
}

func postFile(r Request) *Response {
	if strings.HasPrefix(r.Path, "/files/") && r.Method == "POST" {
		filePath := strings.TrimPrefix(r.Path, "/files/")
		if strings.Contains(filePath, "..") {
			return &Response{
				HttpVersion:   r.HttpVersion,
				StatusCode:    400,
				StatusMessage: "Bad Request",
				Headers:       Headers{},
				Body:          []byte("Invalid file path"),
			}
		}

		fullPath := path.Join(directoryFlag, filePath)
		os.WriteFile(fullPath, r.Body, 0644)
		return &Response{
			HttpVersion:   r.HttpVersion,
			StatusCode:    201,
			StatusMessage: "Created",
			Headers:       Headers{},
			Body:          []byte{},
		}
	}
	return nil
}
