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
// 	"sync"

// 	"github.com/thebeginner86/hippocampus/internal/security"
// 	"github.com/thebeginner86/hippocampus/resp"
// )

// type ConfigHandler struct {
// 	store map[string]string
// 	mu    sync.RWMutex

// 	security *security.Security
// }

// type ConfigSubCmd string

// var (
// 	GetCmd ConfigSubCmd = "get"
// 	SETCmd ConfigSubCmd = "set"
// )

// func config(args []resp.Value) resp.Value {
// 	// means neither get nor set is being requested
// 	// get requires 2 subcmds, eg: CONFIG get save
// 	// set requires 3 subcmds, eg: CONFIG set save "3600 1 300 100 60 10000"
// 	if len(args) != 2 && len(args) != 3 {
// 		return resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'config' command. Must be 2 argument for GET or 3 arguments for SET"}
// 	}

// 	subcmd := strings.ToLower(args[0].Bulk) // makes case in-sensitve

// 	switch subcmd {
// 	case string(GetCmd):
// 		return getConfig(args)
// 	case string(SETCmd):
// 		return setConfig(args)
// 	}
// 	return resp.Value{Type: "error", String: "Error: Invalid subcommand for 'config' command. Only 'get' or 'set'are supported"}
// }
