package resp

import (
	"io"
	"strconv"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// Marshal() is responsible for delegating the serialization based on value type
func (v Value) Marshal() []byte {
  switch v.Type {
    case "array": 
      return v.marshalArray()
    case "bulk": 
      return v.marshalBulk()
    case "string": 
      return v.marshalString()
    case "null":
      return v.marshalNull()
    case "error":
      return v.marshalError()
    default:
      return []byte{}
  }
}


// marshalString() responsible for serializing string value
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.String...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}


// marshalBulk() responsible for serializing bulk value
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}


// marshalArray() responsible for serializing array value
func (v Value) marshalArray() []byte {
	len := len(v.Array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.Array[i].Marshal()...)
	}

	return bytes
}


// marshalInteger() responsible for serializing integer value
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.String...)
	bytes = append(bytes, '\r', '\n')
	
	return bytes
}


// marshalNull() responsible for serializing null value
func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}