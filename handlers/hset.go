package handlers

import (
	"github.com/thebeginner86/hippocampus/resp"
)

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'hset' command. Must be 3 arguments"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMutex.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = make(map[string]string)
	}
	HSETs[hash][key] = value
	HSETsMutex.Unlock()

	return resp.Value{Type: "string", String: "OK"}
}