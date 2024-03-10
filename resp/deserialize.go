package resp

import (
		"io"
    "bufio"
    "fmt"
    "strings"
    "os"
    "strconv"
)

const (
	STRING  = '+'
	BULK    = '$'
	ARRAY   = '*'
)

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}


func deserialize() {
    input := "" // take from user
    reader := bufio.NewReader(strings.NewReader(input))

    b, _ := reader.ReadByte()

    if b != '$' {
      fmt.Println("Invalid type, expecting bulk strings only")
      os.Exit(1)
    }

    size, _ := reader.ReadByte()

    strSize, _ := strconv.ParseInt(string(size), 10, 64)

    // consume /r/n
    reader.ReadByte()
    reader.ReadByte()

    name := make([]byte, strSize)
    reader.Read(name)

    fmt.Println(string(name))
}