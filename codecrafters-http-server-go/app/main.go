package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

var gzipEncodingMiddleware = GzipEncodingMiddleware{}

func main() {
	flag.StringVar(&directoryFlag, "directory", ".", "Directory to serve files from")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		err := serve(l)
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			os.Exit(1)
		}
	}
}

func serve(l net.Listener) error {
	conn, err := l.Accept()
	if err != nil {
		return fmt.Errorf("error accepting connection: %w", err)
	}

	go func() {
		defer conn.Close()

		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			fmt.Println("error reading from connection:", err)
		}
		b = b[:n]
		request, err := parseRequest(string(b))
		if err != nil {
			fmt.Println("error parsing request:", err)
		}

		// response := handleRequest(*request)
		response := gzipEncodingMiddleware.Dispatch(*request, handleRequest)
		fmt.Printf("%v \"%s %s %s\" %d %s\n", now(), request.Method, request.Path, request.HttpVersion, response.StatusCode, response.StatusMessage)

		_, err = conn.Write([]byte(response.String()))
		if err != nil {
			fmt.Println("error writing response:", err)
		}
	}()
	return nil
}
