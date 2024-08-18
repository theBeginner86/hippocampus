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

	"github.com/thebeginner86/hippocampus/internal/handlers"
	"github.com/thebeginner86/hippocampus/resp"
)

var (
	databaseFile = "database.aof"
)

func main() {
	fmt.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	handlerInst, err := handlers.NewHandler(databaseFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer handlerInst.AofH.Close()

	err = handlerInst.AofH.Read(func(value *resp.Value) {
		resp := handlerInst.ExecuteCmd(value, true)
		if resp != nil && resp.Type == "error" {
			return
		}
	})

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
		req, err := respClient.Read()
		if err != nil {
			return
		}

		// this validaiton ensures that the request includes an array and not empty
		if req.Type != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(req.Array) == 0 {
			fmt.Println("Invalid request, empty array")
			continue
		}

		writer := resp.NewWriter(conn)

		resp := handlerInst.ExecuteCmd(req, false)
		writer.Write(resp)
	}
}
