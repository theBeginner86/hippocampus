package config

// import (
// 	"strings"

// 	"github.com/thebeginner86/hippocampus/resp"
// )

// func (handler *ConfigHandler) GET(args []resp.Value) resp.Value {
// 	entity := strings.ToLower(args[1].Bulk) // makes case in-sensitve
// 	if entity != "save" && entity != "appendonly" {
// 		return resp.Value{Type: "error", String: "Error: Invalid configuration entity. Only 'save' or 'appendonly' are supported"}
// 	}

// 	handler.mu.RLock()
// 	configVal, ok := handler.store[entity]
// 	handler.mu.RUnlock()

// 	if !ok {
// 		return resp.Value{Type: "null"}
// 	}

// 	vals := []resp.Value{
// 		{Type: "bulk", Bulk: args[1].Bulk},
// 		{Type: "bulk", Bulk: configVal},
// 	}
// 	return resp.Value{Type: "array", Array: vals}

// }