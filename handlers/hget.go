package handlers

import (
	"github.com/thebeginner86/hippocampus/resp"
)

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'hget' command. Must be 2 argumentor"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMutex.RLock()
	value, ok := HSETs[hash][key]
	HSETsMutex.RUnlock()

	if !ok {
		return resp.Value{Type: "null"}
	}
	
	return resp.Value{Type: "bulk", Bulk: value}
}