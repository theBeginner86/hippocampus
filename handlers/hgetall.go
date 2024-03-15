package handlers

import (
	"github.com/thebeginner86/hippocampus/resp"
)

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'hgetall' command. Must be 1 argumentor"}
	}

	hash := args[0].Bulk

	HSETsMutex.RLock()
	hset, ok := HSETs[hash]
	HSETsMutex.RUnlock()

	if !ok {
		return resp.Value{Type: "null"}
	}

	vals := make([]resp.Value, 0)
	for key, val := range hset {
		vals = append(vals, resp.Value{Type: "bulk", Bulk: key})
		vals = append(vals, resp.Value{Type: "bulk", Bulk: val})
	}

	return resp.Value{Type: "array", Array: vals}
}