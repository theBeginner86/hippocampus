package handlers

import (
	"github.com/thebeginner86/hippocampus/resp"
)

func ping(args []resp.Value) resp.Value {
	if (len(args) == 0) {
		return resp.Value{Type: "string", String: "PONG"}
	}
	
	return resp.Value{Type: "string", String: args[0].Bulk}
}

