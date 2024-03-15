package handlers

import (
	"sync"

	"github.com/thebeginner86/hippocampus/resp"
)

// this hashmap is used to refer the key-value pairs
var SETs = map[string]string{}
// ReadWrite Mutext is necessary:
// 1. This system is responsible for handling multiple requests concurrectly
// 2. To ensure same key is not modified by multiple threads at the same time. That is, To ensure mutual exclusion
var SETsMutex = sync.RWMutex{}

var Handlers = map[string]func([]resp.Value) resp.Value {
	"PING": ping,
	"GET": get,
	"SET": set,
}

