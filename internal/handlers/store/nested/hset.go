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

package nested

import (
	"github.com/thebeginner86/hippocampus/resp"
)

type HSetCmd struct {
	Name string

	*NestedStoreHandler
}

func NewHSetCmd(handler *NestedStoreHandler) *HSetCmd {
	return &HSetCmd{
		Name: "HSET",
		NestedStoreHandler: handler,
	}
}


func (handler *HSetCmd) Handle(req *resp.Value, skp bool) (*resp.Value) {
	if skp {
		return handler.run(req.Array[1:])
	}
	
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

func (handler *HSetCmd) preProcess(req *resp.Value) *resp.Value {
	args := req.Array[1:]
	if len(args) != 3 {
		return &resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'hset' command. Must be 3 arguments"}
	}

	value := args[2].Bulk
	encrptedVal, err := handler.securityH.Encrypter.Encrypt(value)
	if err != nil {
		return &resp.Value{Type: "error", String: "Error: " + err.Error()}
	}

	args[2].Bulk = encrptedVal

	return nil
}

func (handler *HSetCmd) run(args []resp.Value) *resp.Value {

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk
	
	handler.mu.Lock()
	if _, ok := handler.store[hash]; !ok {
		handler.store[hash] = make(map[string]string)
	}
	handler.store[hash][key] = value
	handler.mu.Unlock()

	return nil
}

func (handler *HSetCmd) postProcess(req *resp.Value) *resp.Value {
	err := handler.aofH.Write(req)
	if err != nil {
		return &resp.Value{Type: "error", String: "Error: " + err.Error()}
	}

	return nil
}

