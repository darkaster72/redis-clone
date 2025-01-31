package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, noOfBytes int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		noOfBytes++
		line = append(line, b)

		if len(line) >= 2 && line[len(line)-2] == '\r' && line[len(line)-1] == '\n' {
			break
		}
	}

	return line[:len(line)-2], noOfBytes, nil
}

func (r *Resp) readInteger() (x int, noOfBytes int, err error) {
	line, noOfBytes, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(i64), noOfBytes, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case BULK:
		return r.readBulk()

	case ARRAY:
		return r.readArray()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, fmt.Errorf("invalid type")
	}
}

func (r *Resp) readArray() (Value, error) {
	v := Value{typ: "array"}

	length, _, err := r.readInteger()

	if err != nil {
		return v, err
	}

	v.array = make([]Value, 0)

	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array = append(v.array, val)
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{typ: "bulk"}

	length, _, err := r.readInteger()

	if err != nil {
		return v, err
	}

	bulk := make([]byte, length)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}
