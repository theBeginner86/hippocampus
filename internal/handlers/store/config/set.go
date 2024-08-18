package config

// import (
// 	"strings"

// 	"github.com/thebeginner86/hippocampus/resp"
// )

// func (handler *ConfigHandler) setConfig(args []resp.Value) resp.Value {
// 	entity := strings.ToLower(args[1].Bulk) // makes case in-sensitve
// 	if entity != "save" && entity != "appendonly" {
// 		return resp.Value{Type: "error", String: "Error: Invalid configuration entity. Only 'save' or 'appendonly' are supported"}
// 	}

// 	value := args[2].Bulk

// 	handler.mu.RLock()
// 	handler.store[entity] = value
// 	handler.mu.RUnlock()

// 	return resp.Value{Type: "string", String: "OK"}
// }
