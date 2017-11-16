// Code generated by trunnel. DO NOT EDIT.

package unionlo

import (
	"encoding/binary"
	"errors"
)

type Unlo struct {
	Tag         uint8
	X           uint8
	Y           []uint8
	Byte        uint8
	Z           []uint16
	Leftoverlen uint8
	Leftovers   []uint8
}

func (u *Unlo) Parse(data []byte) ([]byte, error) {
	{
		if len(data) < 1 {
			return nil, errors.New("data too short")
		}
		u.Tag = data[0]
		data = data[1:]
	}
	{
		if len(data) < 8 {
			return nil, errors.New("data too short")
		}
		restore := data[len(data)-8:]
		data = data[:len(data)-8]
		switch {
		case u.Tag == 1:
			{
				if len(data) < 1 {
					return nil, errors.New("data too short")
				}
				u.X = data[0]
				data = data[1:]
			}
		case u.Tag == 2:
			{
				u.Y = make([]uint8, 0)
				for len(data) > 0 {
					var t uint8
					if len(data) < 1 {
						return nil, errors.New("data too short")
					}
					t = data[0]
					data = data[1:]
					u.Y = append(u.Y, t)
				}
			}
		case u.Tag == 4:
			{
				if len(data) < 1 {
					return nil, errors.New("data too short")
				}
				u.Byte = data[0]
				data = data[1:]
			}
			{
				u.Z = make([]uint16, 0)
				for len(data) > 0 {
					var t uint16
					if len(data) < 2 {
						return nil, errors.New("data too short")
					}
					t = binary.BigEndian.Uint16(data)
					data = data[2:]
					u.Z = append(u.Z, t)
				}
			}
		}
		if len(data) > 0 {
			return nil, errors.New("trailing data disallowed")
		}
		data = restore
	}
	{
		if len(data) < 1 {
			return nil, errors.New("data too short")
		}
		u.Leftoverlen = data[0]
		data = data[1:]
	}
	{
		u.Leftovers = make([]uint8, int(u.Leftoverlen))
		for i := 0; i < int(u.Leftoverlen); i++ {
			if len(data) < 1 {
				return nil, errors.New("data too short")
			}
			u.Leftovers[i] = data[0]
			data = data[1:]
		}
	}
	return data, nil
}
