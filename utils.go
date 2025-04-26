package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

var (
	internalAddress = "localhost:1238"
)

func receiveOnceTCPServer(size int, readFunc func(io.ReadCloser) (int, error)) error {
	ln, err := net.Listen("tcp", internalAddress)
	if err != nil {
		fmt.Println("Server failed to start:", err)
		return err
	}
	defer ln.Close()
	fmt.Println("Server listening on", internalAddress)

	// give the server a sec to start
	time.Sleep(1 * time.Second)

	go sendTCPMessage(size)

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Failed to accept connection:", err)
		return err
	}
	defer conn.Close()

	n, err := readFunc(conn)
	fmt.Printf("Server received %d bytes\n", n)
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

func sendTCPMessage(size int) {
	fmt.Println("Sending message to", internalAddress)
	conn, err := net.Dial("tcp", internalAddress)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	message := strings.Repeat("A", size)

	n, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}
	if n != size {
		fmt.Printf("Failed to send message: sent insufficient size=%d expectedSize=%d\n", n, size)
		return
	}

	fmt.Printf("Client finished sending %d bytes\n", size)
}
