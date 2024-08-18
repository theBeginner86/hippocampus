package std

import (
	"github.com/thebeginner86/hippocampus/resp"
)

type GetAllCmd struct {
	Name string
	*StdStoreHandler
}


func NewGetAllCmd(handler *StdStoreHandler) *GetAllCmd {
	return &GetAllCmd{
		Name: "GETALL",
		StdStoreHandler: handler,
	}
}

func (handler *GetAllCmd) Handle(req *resp.Value) *resp.Value {
	res := handler.preProcess(req)
	if res != nil && res.Type == "error" {
		return res
	}

	res = handler.run(req.Array[1:])
	
	return handler.postProcess(res)
}

func (handler *GetAllCmd) preProcess(req *resp.Value) *resp.Value {
	args := req.Array[1:]
	if len(args) != 1 {
		return &resp.Value{Type: "error", String: "Error: Invalid number of arguments for 'getall' command. Must be 1 argument"}
	}

	return nil
}

// Note: This is a custom handler for GETALL command. Not seen is current redis version
func (handler *GetAllCmd) run(_ []resp.Value) *resp.Value {
	handler.mu.RLock()
	defer handler.mu.RUnlock()

	vals := make([]resp.Value, 0)
	for key, value := range handler.store {
		vals = append(vals, resp.Value{Type: "bulk", Bulk: key})
		vals = append(vals, resp.Value{Type: "bulk", Bulk: value})
	}

	return &resp.Value{Type: "array", Array: vals}
}

func (handler *GetAllCmd) postProcess(req *resp.Value) *resp.Value {
	decryptedVals := make([]resp.Value, 0)
	for idx, val := range req.Array {
		if idx%2 == 0 {
			decryptedVal, err := handler.securityH.Decrypter.Decrypt(val.Bulk)
			
			// TODO: Handle error less promptly? allow some fields to be returned?
			if err != nil {
				return &resp.Value{Type: "error", String: "Error: " + err.Error()}
			}
			decryptedVals = append(decryptedVals, resp.Value{Type: "bulk", Bulk: decryptedVal})
		} else {
			decryptedVals = append(decryptedVals, val)
		}
	}

	return &resp.Value{Type: "array", Array: decryptedVals}
}