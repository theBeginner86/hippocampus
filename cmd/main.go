//  Copyright 2024 Pranav Singh

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/thebeginner86/hippocampus/internal/handlers"
	"github.com/thebeginner86/hippocampus/persistance/aof"
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

	// TODO: change it to hidden folder
	aofH, err := aof.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aofH.Close()

	aofH.Read(func(value resp.Value) {
		cmd := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := handlers.Handlers[cmd]
		if !ok {
			fmt.Println("Invalid command: ", cmd)
			return
		}

		handler(args)
	})

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

		// fetches the first element of the array and convert to uppercase
		// this ensures that the cmds matches properly with our defined standards
		command := strings.ToUpper(value.Array[0].Bulk)

		// rest of the elements would be arguments of the redis cmds
		args := value.Array[1:]

		writer := resp.NewWriter(conn)

		handler, ok := handlers.Handlers[command]
		// if command == "COMMAND" {continue}  
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(resp.Value{Type: "string", String: ""})
			continue
		}

		if command == "SET" ||
			command == "HSET" ||
			(command == "CONFIG" && len(args) > 0 &&
				strings.ToUpper(args[0].Bulk) == "SET") {
			aofH.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
