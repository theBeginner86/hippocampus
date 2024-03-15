package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/thebeginner86/hippocampus/handlers"
	"github.com/thebeginner86/hippocampus/resp"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server of type tcp listener
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
			return
		}

		// this validaiton ensures that the request includes an array and not empty
		if value.Type != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, empty array")
			continue
		}

		// we fetch the first element of the array and convert to uppercase
		// this ensures that the cmds matches properly with our defined standards
		command := strings.ToUpper(value.Array[0].Bulk)

		// rest of the elements would be arguments of the redis cmds
		args := value.Array[1:]

		writer := resp.NewWriter(conn)

		handler, ok := handlers.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(resp.Value{Type: "string", String: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
