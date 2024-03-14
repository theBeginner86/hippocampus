package main

import (
	"fmt"
	"net"

	"github.com/thebeginner86/hippocampus/resp"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	// this is a tcp listener
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		respClient := resp.NewResp(conn)
		value, err := respClient.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(value)

		// ignore request and return back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
}