package resp

import (
	"fmt"
	"strconv"
)

// TODO: deprecate this and use in-build function from bufio
// readLine() is a helper function to read a line from the reader. This continuosly reads 
// untill it finds \r\n
func (r *Resp) readLine() (line []byte, length int, err error){
	for{
		byt, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		length++
		line = append(line, byt)
		// its should be greater than 2 because first 2 bytes are type and length
		// this would be true when we are adding \n to line, that is, the last byte
		// at that point second last character would be \r
		if len(line) >= 2 && line[len(line) - 2] == '\r' { 
			break
		}
	}

	// return line withouth trailing \r\n as this would be added later by the caller
	return line[:len(line) - 2], length, nil
}


func (r *Resp) readInteger() (integer int, length int, err error){
	line, n, err := r.readLine();
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, length, err
	}
	return int(i64), n, nil
}


func (r *Resp) readBulk() (Value, error){
	val := Value{
		Type: "bulk",
	}
	len, _, err := r.readInteger()
	if err != nil {
		return val, err
	}
	bulk := make([]byte, len)
	r.reader.Read(bulk)
	val.Bulk = string(bulk)
	
	r.readLine() // read trailing CRLF

	return val, nil
}

// readArray() is the handler for reading an array from the reader
func (r *Resp) readArray() (Value, error){
	val := Value{
		Type: "array",
	}

	len, _, err := r.readInteger()
	if err != nil {
		return val, err
	}

	val.Array = make([]Value, len)
	for i:=0; i<len; i++ {
		val, err := r.Read()
		if err != nil {
			return val, err
		}
		val.Array = append(val.Array, val)
	}
	return val, nil
}


// Read() is an entry point to read a RESP value
// based on the type it delegates the requests to the requested handler
func (r *Resp) Read() (Value, error) {
	typ, err := r.reader.ReadByte()
	if err !=nil {
		return Value{}, err
	}

	switch typ {
		case ARRAY:
			return r.readArray()
		case BULK:
			return r.readBulk()
		default:
			fmt.Printf("Unkown type: %v", string(typ))
			return Value{}, nil
	}

}