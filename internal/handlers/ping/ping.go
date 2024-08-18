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

package ping

import (
	"github.com/thebeginner86/hippocampus/resp"
)

type PingCmd struct {
	Name string
}

func NewPingCmd() *PingCmd {
	return &PingCmd{
		Name: "PING",
	}
}

func (handler *PingCmd) Handle(req *resp.Value) *resp.Value {
	res := handler.preProcess(req)
	if res != nil && res.Type != "error" {
		return res
	}
	res = handler.run(req.Array[1:])
	if res != nil && res.Type == "error" {
		return res
	}
	return handler.postProcess(res)
}

func (handler *PingCmd) preProcess(req *resp.Value) *resp.Value {
	if len(req.Array) > 1 {
		return &resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'ping' command. Must be 0 or 1 argument"}
	}

	return nil
}

func (handler *PingCmd) run(args []resp.Value) *resp.Value {
	if len(args) == 0 {
		return &resp.Value{Type: "string", String: "PONG"}
	}

	return &resp.Value{Type: "string", String: args[0].Bulk}
}

func (handler *PingCmd) postProcess(req *resp.Value) *resp.Value {
	return req
}
