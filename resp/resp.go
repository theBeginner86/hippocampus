package resp

import (
		"io"
    "bufio"
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
  typ string
  str string
  num int
  bulk string
  array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}