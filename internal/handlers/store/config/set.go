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

package config

// import (
// 	"strings"

// 	"github.com/thebeginner86/hippocampus/resp"
// )

// func (handler *ConfigHandler) setConfig(args []resp.Value) resp.Value {
// 	entity := strings.ToLower(args[1].Bulk) // makes case in-sensitve
// 	if entity != "save" && entity != "appendonly" {
// 		return resp.Value{Type: "error", String: "Error: Invalid configuration entity. Only 'save' or 'appendonly' are supported"}
// 	}

// 	value := args[2].Bulk

// 	handler.mu.RLock()
// 	handler.store[entity] = value
// 	handler.mu.RUnlock()

// 	return resp.Value{Type: "string", String: "OK"}
// }
