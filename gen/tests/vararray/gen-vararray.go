// Code generated by trunnel. DO NOT EDIT.

package vararray

import (
	"encoding/binary"
	"errors"
)

type VarArray struct {
	NWords uint16
	Words  []uint32
}

func (v *VarArray) Parse(data []byte) ([]byte, error) {
	{
		if len(data) < 2 {
			return nil, errors.New("data too short")
		}
		v.NWords = binary.BigEndian.Uint16(data)
		data = data[2:]
	}
	{
		v.Words = make([]uint32, int(v.NWords))
		for i := 0; i < int(v.NWords); i++ {
			if len(data) < 4 {
				return nil, errors.New("data too short")
			}
			v.Words[i] = binary.BigEndian.Uint32(data)
			data = data[4:]
		}
	}
	return data, nil
}