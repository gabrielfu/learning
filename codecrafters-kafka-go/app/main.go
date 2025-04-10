package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	err := startServer()
	if err != nil {
		fmt.Println("Error starting server: ", err.Error())
		os.Exit(1)
	}
}

func startServer() error {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		return fmt.Errorf("failed to bind to port 9092")
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go func() {
			err := handleConnection(conn)
			if err != nil && err != io.EOF {
				fmt.Println("Error handling connection: ", err.Error())
			}
		}()
	}
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()
	for {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		if err != nil {
			return err
		}
		b = b[:n]

		var request Request
		err = request.UnmarshalBinary(b)
		if err != nil {
			return err
		}

		response := handleRequest(request)
		rb, _ := response.MarshalBinary()
		conn.Write(rb)
	}
}
