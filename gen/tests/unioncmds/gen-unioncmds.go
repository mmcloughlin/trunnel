// Code generated by trunnel. DO NOT EDIT.

package unioncmds

import (
	"encoding/binary"
	"errors"
)

type UnionCmds struct {
	Tag uint8
	X   [2]uint32
	Y   uint32
}

func (u *UnionCmds) Parse(data []byte) ([]byte, error) {
	cur := data
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		u.Tag = cur[0]
		cur = cur[1:]
	}
	{
		switch {
		case u.Tag == 1:
		case u.Tag == 2:
			{
				return nil, errors.New("disallowed case")
			}
		default:
			{
				for idx := 0; idx < 2; idx++ {
					if len(cur) < 4 {
						return nil, errors.New("data too short")
					}
					u.X[idx] = binary.BigEndian.Uint32(cur)
					cur = cur[4:]
				}
			}
		}
	}
	{
		if len(cur) < 4 {
			return nil, errors.New("data too short")
		}
		u.Y = binary.BigEndian.Uint32(cur)
		cur = cur[4:]
	}
	return cur, nil
}