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
	cur := data
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		u.Tag = cur[0]
		cur = cur[1:]
	}
	{
		if len(cur) < 8 {
			return nil, errors.New("data too short")
		}
		restore := cur[len(cur)-8:]
		cur = cur[:len(cur)-8]
		switch {
		case u.Tag == 1:
			{
				if len(cur) < 1 {
					return nil, errors.New("data too short")
				}
				u.X = cur[0]
				cur = cur[1:]
			}
		case u.Tag == 2:
			{
				u.Y = make([]uint8, 0)
				for len(cur) > 0 {
					var tmp uint8
					if len(cur) < 1 {
						return nil, errors.New("data too short")
					}
					tmp = cur[0]
					cur = cur[1:]
					u.Y = append(u.Y, tmp)
				}
			}
		case u.Tag == 4:
			{
				if len(cur) < 1 {
					return nil, errors.New("data too short")
				}
				u.Byte = cur[0]
				cur = cur[1:]
			}
			{
				u.Z = make([]uint16, 0)
				for len(cur) > 0 {
					var tmp uint16
					if len(cur) < 2 {
						return nil, errors.New("data too short")
					}
					tmp = binary.BigEndian.Uint16(cur)
					cur = cur[2:]
					u.Z = append(u.Z, tmp)
				}
			}
		}
		if len(cur) > 0 {
			return nil, errors.New("trailing data disallowed")
		}
		cur = restore
	}
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		u.Leftoverlen = cur[0]
		cur = cur[1:]
	}
	{
		u.Leftovers = make([]uint8, int(u.Leftoverlen))
		for idx := 0; idx < int(u.Leftoverlen); idx++ {
			if len(cur) < 1 {
				return nil, errors.New("data too short")
			}
			u.Leftovers[idx] = cur[0]
			cur = cur[1:]
		}
	}
	return cur, nil
}
