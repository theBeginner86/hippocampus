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

package handlers

import (
	"github.com/thebeginner86/hippocampus/resp"
)

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'hgetall' command. Must be 1 argumentor"}
	}

	hash := args[0].Bulk

	HSETsMutex.RLock()
	hset, ok := HSETs[hash]
	HSETsMutex.RUnlock()

	if !ok {
		return resp.Value{Type: "null"}
	}

	vals := make([]resp.Value, 0)
	for key, val := range hset {
		vals = append(vals, resp.Value{Type: "bulk", Bulk: key})
		vals = append(vals, resp.Value{Type: "bulk", Bulk: val})
	}

	return resp.Value{Type: "array", Array: vals}
}
