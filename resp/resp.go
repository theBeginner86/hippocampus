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

package resp

import (
	"bufio"
	"io"
)

// RESP protocol supported types
const (
	STRING  = '+'
	BULK    = '$'
	ARRAY   = '*'
	INTEGER = ":"
	ERROR   = '-'
)

// represents a RESP value
// used for serializing/deserializing RESP cmds
type Value struct {
	Type   string
	String string
	Number int
	Bulk   string
	Array  []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}
