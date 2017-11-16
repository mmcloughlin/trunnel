// Code generated by trunnel. DO NOT EDIT.

package unionlen

import (
	"encoding/binary"
	"errors"
)

type UnionWithLen struct {
	Tag                uint16
	UnionLen           uint16
	R                  uint8
	G                  uint8
	B                  uint8
	Year               uint16
	Month              uint8
	Day                uint8
	Unparseable        []uint8
	RightAfterTheUnion uint16
}

func (u *UnionWithLen) Parse(data []byte) ([]byte, error) {
	{
		if len(data) < 2 {
			return nil, errors.New("data too short")
		}
		u.Tag = binary.BigEndian.Uint16(data)
		data = data[2:]
	}
	{
		if len(data) < 2 {
			return nil, errors.New("data too short")
		}
		u.UnionLen = binary.BigEndian.Uint16(data)
		data = data[2:]
	}
	{
		if len(data) < int(u.UnionLen) {
			return nil, errors.New("data too short")
		}
		restore := data[int(u.UnionLen):]
		data = data[:int(u.UnionLen)]
		switch {
		case u.Tag == 1:
			{
				if len(data) < 1 {
					return nil, errors.New("data too short")
				}
				u.R = data[0]
				data = data[1:]
			}
			{
				if len(data) < 1 {
					return nil, errors.New("data too short")
				}
				u.G = data[0]
				data = data[1:]
			}
			{
				if len(data) < 1 {
					return nil, errors.New("data too short")
				}
				u.B = data[0]
				data = data[1:]
			}
		case u.Tag == 2:
			{
				if len(data) < 2 {
					return nil, errors.New("data too short")
				}
				u.Year = binary.BigEndian.Uint16(data)
				data = data[2:]
			}
			{
				if len(data) < 1 {
					return nil, errors.New("data too short")
				}
				u.Month = data[0]
				data = data[1:]
			}
			{
				if len(data) < 1 {
					return nil, errors.New("data too short")
				}
				u.Day = data[0]
				data = data[1:]
			}
			{
				data = []byte{}
			}
		default:
			{
				u.Unparseable = make([]uint8, 0)
				for len(data) > 0 {
					var t uint8
					if len(data) < 1 {
						return nil, errors.New("data too short")
					}
					t = data[0]
					data = data[1:]
					u.Unparseable = append(u.Unparseable, t)
				}
			}
		}
		if len(data) > 0 {
			return nil, errors.New("trailing data disallowed")
		}
		data = restore
	}
	{
		if len(data) < 2 {
			return nil, errors.New("data too short")
		}
		u.RightAfterTheUnion = binary.BigEndian.Uint16(data)
		data = data[2:]
	}
	return data, nil
}
