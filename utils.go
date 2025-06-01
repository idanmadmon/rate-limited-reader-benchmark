package main

import (
	"fmt"
	"io"
	"net"
)

var (
	serverAddress = "localhost:1238"
)

type readFunc func(io.ReadCloser) (int, error)
type writeFunc func(io.Writer) (int, error)

func receiveOnceTCPServer(rf readFunc) (int, error) {
	ln, err := net.Listen("tcp", serverAddress)
	if err != nil {
		fmt.Println("Server failed to start:", err)
		return 0, err
	}
	defer ln.Close()
	fmt.Println("Server listening on", serverAddress)

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Failed to accept connection:", err)
		return 0, err
	}
	defer conn.Close()

	n, err := rf(conn)
	fmt.Printf("Server received %d bytes\n", n)
	return n, err
}

func sendTCPMessage(wf writeFunc) (int, error) {
	fmt.Println("Sending message to", serverAddress)
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return 0, err
	}
	defer conn.Close()

	n, err := wf(conn)
	fmt.Printf("Client sent %d bytes\n", n)
	return n, err
}
