package handlers

import (
	"github.com/thebeginner86/hippocampus/resp"
)

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'get' command. Must be 1 argument"}
	}

	key := args[0].Bulk

	SETsMutex.RLock()
	value, ok := SETs[key]
	SETsMutex.RUnlock()

	if !ok {
		return resp.Value{Type: "null"}
	}

	return resp.Value{Type: "bulk", Bulk: value}
}