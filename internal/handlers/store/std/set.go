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

type SetCmd struct {
	Name string

	*StdStoreHandler
}

func NewSetCmd(handler *StdStoreHandler) *SetCmd {
	return &SetCmd{
		Name: "SET",
		StdStoreHandler: handler,
	}
}

func (handler *SetCmd) Handle(req *resp.Value) (*resp.Value) {
	res := handler.preProcess(req)	
	if res != nil && res.Type != "error" {
		return res
	}

	handler.run(req.Array[1:])

	res = handler.postProcess(req)
	if res != nil && res.Type != "error" {
		return res
	}

	return &resp.Value{Type: "string", String: "OK"}
}

func (handler *SetCmd) preProcess(req *resp.Value) *resp.Value {
	args := req.Array[1:]
	if len(args) != 2 {
		return &resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'set' command. Must be 2 arguments"}
	}

	value := args[1].Bulk
	encrptedVal, err := handler.securityH.Encrypter.Encrypt(value)
	if err != nil {
		return &resp.Value{Type: "error", String: "Error: " + err.Error()}
	}

	args[1].Bulk = encrptedVal

	return nil
}

func (handler *SetCmd) run(args []resp.Value) *resp.Value{
	key := args[0].Bulk
	value := args[1].Bulk

	handler.mu.Lock()
	handler.store[key] = value
	handler.mu.Unlock()

	return nil
}

func (handler *SetCmd) postProcess(req *resp.Value) (*resp.Value) {
	err := handler.aofH.Write(req)
	if err != nil {
		return &resp.Value{Type: "error", String: "Error: " + err.Error()}
	}
	return nil
}
