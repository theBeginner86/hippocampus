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