//  Copyright 2024 Pranav Singh

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//      http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package handlers

import (
	"strings"

	"github.com/thebeginner86/hippocampus/internal/handlers/store/nested"
	"github.com/thebeginner86/hippocampus/internal/handlers/store/std"
	"github.com/thebeginner86/hippocampus/internal/handlers/ping"
	"github.com/thebeginner86/hippocampus/internal/security"
	"github.com/thebeginner86/hippocampus/persistance/aof"
	"github.com/thebeginner86/hippocampus/resp"
)

type Cmds string

const (
	SET Cmds = "SET"
	GET Cmds = "GET"
	GETALL Cmds = "GETALL"
	HSET Cmds = "HSET"
	HGET Cmds = "HGET"
	HGETALL Cmds = "HGETALL"
	CONFIG Cmds = "CONFIG"
  PING Cmds = "PING"
)

type IStoreHandler interface {
	Handle(req *resp.Value, skip bool) *resp.Value // TODO: skp should be removed
	preProcess(req *resp.Value) *resp.Value 
	run(args []resp.Value) *resp.Value
	postProcess(req *resp.Value) *resp.Value
}

type CmdsHandler struct {
	SetCmdH *std.SetCmd
	GetCmdH *std.GetCmd
	GetAllCmdH *std.GetAllCmd

	HSetCmdH *nested.HSetCmd
	HGetCmdH *nested.HGetCmd
	HGetAllCmdH *nested.HGetAllCmd

	PingCmdH *ping.PingCmd
}

func newCmdHandler(secH *security.Security, aofH *aof.Aof) *CmdsHandler {
	stdH := std.NewStdStoreHandler(secH, aofH)
	nestedH := nested.NewNestedStoreHandler(secH, aofH)
	return &CmdsHandler{
		SetCmdH: std.NewSetCmd(stdH),
		GetCmdH: std.NewGetCmd(stdH),
		GetAllCmdH: std.NewGetAllCmd(stdH),
		HGetCmdH: nested.NewHGetCmd(nestedH),
		HSetCmdH: nested.NewHSetCmd(nestedH),
		HGetAllCmdH: nested.NewHGetAllCmd(nestedH),
		PingCmdH: ping.NewPingCmd(),
	}
}

type Handler struct {
	CmdHander *CmdsHandler
	AofH *aof.Aof
}

var byts = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
const privateKey string = "abc&1*~#^2^#s0^=)^^7%b34"

func NewHandler(aofFile string) (*Handler, error) {
  aofH, err := aof.NewAof(aofFile)
	if err != nil {
		return nil, err
	}	
	secH := security.NewSecurity(privateKey, byts)

	return &Handler{
		CmdHander: newCmdHandler(secH, aofH),
		AofH: aofH,
	}, nil
}

func (handler *Handler) ExecuteCmd(req *resp.Value, skp bool) *resp.Value {	
	// fetches the first element of the array and convert to uppercase
	// this ensures that the cmds matches properly with our defined standards
	cmd := strings.ToUpper(req.Array[0].Bulk)
	req.Array[0].Bulk = cmd // TODO: please, fix this oddity 

	switch cmd {
	case string(SET):
		return handler.CmdHander.SetCmdH.Handle(req, skp)
	case string(GET):
		return handler.CmdHander.GetCmdH.Handle(req)
	case string(GETALL):
		return handler.CmdHander.GetAllCmdH.Handle(req)
	case string(HSET):
		return handler.CmdHander.HSetCmdH.Handle(req, skp)
	case string(HGET):
		return handler.CmdHander.HGetCmdH.Handle(req)
	case string(HGETALL):
		return handler.CmdHander.HGetAllCmdH.Handle(req)
	// case string(CONFIG):
	// 	return handler.CmdHander.StdStoreHandler.CONFIG(args)
	case string(PING):
		return handler.CmdHander.PingCmdH.Handle(req)
	default:
		return &resp.Value{Type: "error", Bulk: "Invalid command"}
	}
}
