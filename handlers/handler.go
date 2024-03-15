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
var Handlers = map[string]func([]resp.Value) resp.Value {
	"PING": ping,
	"GET": get,
	"SET": set,
	"HGET": hget,
	"HSET": hset,
}

