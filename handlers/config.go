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
	"strings"
	"fmt"

	"github.com/thebeginner86/hippocampus/resp"
)

func config(args []resp.Value) resp.Value {
	// means neither get nor set is being requested
	// get requires 2 subcmds, eg: CONFIG get save
	// set requires 3 subcmds, eg: CONFIG set save "3600 1 300 100 60 10000"
	if len(args) != 2 && len(args) != 3 {
		return resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'config' command. Must be 2 argument for GET or 3 arguments for SET"}
	}

	subcmd := strings.ToLower(args[0].Bulk) // makes case in-sensitve

	switch subcmd {
		case "get": return getConfig(args);
		case "set": return setConfig(args);
	}
	return resp.Value{Type: "error", String: "Error: Invalid subcommand for 'config' command. Only 'get' or 'set'are supported"}
}

// getConfig gets the configuration value for the given entity
func getConfig(args []resp.Value) resp.Value {
	entity := strings.ToLower(args[1].Bulk) // makes case in-sensitve
	if entity != "save" && entity != "appendonly" {
		return resp.Value{Type: "error", String: "Error: Invalid configuration entity. Only 'save' or 'appendonly' are supported"}
	}

	CONFIGMutex.RLock()
	configVal, ok := CONFIG[entity]
	CONFIGMutex.RUnlock()

	if !ok {
		return resp.Value{Type: "null"}
	}

	vals := []resp.Value{
		{Type: "bulk", Bulk: args[1].Bulk},
		{Type: "bulk", Bulk: configVal},
	}
	return resp.Value{Type: "array", Array: vals}

}

// setConfig sets the configuration value for the given entity
func setConfig(args []resp.Value) resp.Value {
	entity := strings.ToLower(args[1].Bulk) // makes case in-sensitve
	if entity != "save" && entity != "appendonly" {
		return resp.Value{Type: "error", String: "Error: Invalid configuration entity. Only 'save' or 'appendonly' are supported"}
	}

	value := args[2].Bulk

	CONFIGMutex.RLock()	
	CONFIG[entity] = value
	CONFIGMutex.RUnlock()

	fmt.Printf("CONFIG: %v\n", CONFIG)

	return resp.Value{Type: "string", String: "OK"}
}