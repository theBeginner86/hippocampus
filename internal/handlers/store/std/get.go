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

package std

import (
	"github.com/thebeginner86/hippocampus/resp"
)

type GetCmd struct {
	Name string

	*StdStoreHandler
}

func NewGetCmd(handler *StdStoreHandler) *GetCmd {
	return &GetCmd{
		Name: "GET",
		StdStoreHandler: handler,
	}
}

func (handler *GetCmd) Handle(req *resp.Value) *resp.Value {
	res := handler.preProcess(req)
	if res != nil && res.Type != "error" {
		return res
	}

	res = handler.run(req.Array[1:])
	if res != nil && res.Type == "error" {
		return res
	}

	res = handler.postProcess(res)

	return &resp.Value{Type: "bulk", Bulk: res.Bulk}
}

func (handler *GetCmd) preProcess(req *resp.Value) *resp.Value {
	args := req.Array[1:]
	if len(args) != 1 {
		return &resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'get' command. Must be 1 argument"}
	}

	return nil
}

func (handler *GetCmd) run(args []resp.Value) *resp.Value {
	key := args[0].Bulk

	handler.mu.RLock()
	value, ok := handler.store[key]
	handler.mu.RUnlock()

	if !ok {
		return &resp.Value{Type: "null"}
	}

	return &resp.Value{Type: "bulk", Bulk: value}
}

func (handler *GetCmd) postProcess(req *resp.Value) *resp.Value {
	value, err := handler.securityH.Decrypter.Decrypt(req.Bulk)
	if err != nil {
		return &resp.Value{Type: "error", String: "Error: " + err.Error()}
	}	

	return &resp.Value{Type: "bulk", Bulk: value}
}