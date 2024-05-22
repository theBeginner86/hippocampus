package handlers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/thebeginner86/hippocampus/persistance/aof"
	"github.com/thebeginner86/hippocampus/resp"
)

func BenchmarkHandlers(b *testing.B) {
	aofH, err := aof.NewAof("database-1000.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aofH.Close()



// set := []resp.Value{
// 		{Type: "bulk", Bulk: "key"},
// 		{Type: "bulk", Bulk: "value"},
// 	}

	// get := []resp.Value{
	// 	{Type: "bulk", Bulk: "mykey"},
	// }

	for i := 0; i < b.N; i++ {
		// handler := Handlers["SET"]
		// handler(set)
				aofH.Read(func(value resp.Value) {
			cmd := strings.ToUpper(value.Array[0].Bulk)
			args := value.Array[1:]

			handler, ok := Handlers[cmd]
			if !ok {
				fmt.Println("Invalid command: ", cmd)
				return
			}

			handler(args)
		})
	}
}	