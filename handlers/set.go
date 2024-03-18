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

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'set' command. Must be 2 arguments"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMutex.Lock()
	SETs[key] = value
	SETsMutex.Unlock()

	return resp.Value{Type: "string", String: "OK"}
}
