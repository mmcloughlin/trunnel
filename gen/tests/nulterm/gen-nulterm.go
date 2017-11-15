// Code generated by trunnel. DO NOT EDIT.

package nulterm

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Nulterm struct {
	X uint32
	S string
	Y uint8
}

func (n *Nulterm) Parse(data []byte) ([]byte, error) {
	{
		if len(data) < 4 {
			return nil, errors.New("data too short")
		}
		n.X = binary.BigEndian.Uint32(data)
		data = data[4:]
	}
	{
		i := bytes.IndexByte(data, 0)
		if i < 0 {
			return nil, errors.New("could not parse nul-term string")
		}
		n.S, data = string(data[:i]), data[i+1:]
	}
	{
		if len(data) < 1 {
			return nil, errors.New("data too short")
		}
		n.Y = data[0]
		data = data[1:]
	}
	return data, nil
}