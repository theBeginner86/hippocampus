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

type HGetAllCmd struct {
	Name string
	*NestedStoreHandler
}

func NewHGetAllCmd(handler *NestedStoreHandler) *HGetAllCmd {
	return &HGetAllCmd{
		Name:               "HGETALL",
		NestedStoreHandler: handler,
	}
}

func (handler *HGetAllCmd) Handle(req *resp.Value) *resp.Value {
	res := handler.preProcess(req)
	if res != nil && res.Type == "error" {
		return res
	}

	res = handler.run(req.Array[1:])
	if res != nil && res.Type == "error" {
		return res
	}

	return handler.postProcess(res)
}

func (handler *HGetAllCmd) preProcess(req *resp.Value) *resp.Value {
	args := req.Array[1:]
	if len(args) != 1 {
		return &resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'HGetAllCmd' command. Must be 1 argument"}
	}

	return nil
}

func (handler *HGetAllCmd) run(args []resp.Value) *resp.Value {
	hash := args[0].Bulk

	handler.mu.RLock()
	hset, ok := handler.store[hash]
	handler.mu.RUnlock()

	if !ok {
		return &resp.Value{Type: "null"}
	}

	vals := make([]resp.Value, 0)
	for key, val := range hset {
		vals = append(vals, resp.Value{Type: "bulk", Bulk: key})
		vals = append(vals, resp.Value{Type: "bulk", Bulk: val})
	}

	return &resp.Value{Type: "array", Array: vals}
}

func (handler *HGetAllCmd) postProcess(req *resp.Value) *resp.Value {
	decryptedVals := make([]resp.Value, 0)
	for idx, val := range req.Array {
		if idx%2 == 1 {
			decryptedVal, err := handler.securityH.Decrypter.Decrypt(val.Bulk)
			// TODO: Handle error less promptly? allow some fields to be returned?
			if err != nil {
				return &resp.Value{Type: "error", String: "Error: " + err.Error()}
			}
			decryptedVals = append(decryptedVals, resp.Value{Type: "bulk", Bulk: decryptedVal})
		} else {
			decryptedVals = append(decryptedVals, val)
		}
	}

	return &resp.Value{Type: "array", Array: decryptedVals}
}
