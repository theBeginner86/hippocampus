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
	"sync"

	"github.com/thebeginner86/hippocampus/resp"
)

// Datastructures for SETs
// this hashmap is used to refer the key-value pairs
var SETs = map[string]string{}

// ReadWrite Mutex is necessary:
// 1. This system is responsible for handling multiple requests concurrectly
// 2. To ensure same key is not modified by multiple threads at the same time. That is, to ensure mutual exclusion
var SETsMutex = sync.RWMutex{}

// Datastructures for HSETs
// this nested hashmap is used to refer the key-value pairs
var HSETs = map[string]map[string]string{}

// ReadWrite Mutex is necessary:
// 1. This system is responsible for handling multiple requests concurrectly
// 2. To ensure same key is not modified by multiple threads at the same time. That is, to ensure mutual exclusion
var HSETsMutex = sync.RWMutex{}

// Handlers is a map of command names to their respective functions
var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING":    ping,
	"GET":     get,
	"SET":     set,
	"HGET":    hget,
	"HSET":    hset,
	"HGETALL": hgetall,
}
