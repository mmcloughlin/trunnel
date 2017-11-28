// Code generated by trunnel. DO NOT EDIT.

package unionbasic

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Date struct {
	Year  uint16
	Month uint8
	Day   uint8
}

func (d *Date) Parse(data []byte) ([]byte, error) {
	cur := data
	{
		if len(cur) < 2 {
			return nil, errors.New("data too short")
		}
		d.Year = binary.BigEndian.Uint16(cur)
		cur = cur[2:]
	}
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		d.Month = cur[0]
		cur = cur[1:]
	}
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		d.Day = cur[0]
		cur = cur[1:]
	}
	return cur, nil
}

func ParseDate(data []byte) (*Date, error) {
	d := new(Date)
	_, err := d.Parse(data)
	if err != nil {
		return nil, err
	}
	return d, nil
}

type Basic struct {
	Tag        uint8
	D          *Date
	Num        uint32
	Eightbytes [8]uint8
	String     string
}

func (b *Basic) Parse(data []byte) ([]byte, error) {
	cur := data
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		b.Tag = cur[0]
		if !(b.Tag == 2 || b.Tag == 3 || b.Tag == 4 || b.Tag == 5 || b.Tag == 6) {
			return nil, errors.New("integer constraint violated")
		}
		cur = cur[1:]
	}
	{
		switch {
		case b.Tag == 2:
			{
				var err error
				b.D = new(Date)
				cur, err = b.D.Parse(cur)
				if err != nil {
					return nil, err
				}
			}
		case b.Tag == 3:
			{
				if len(cur) < 4 {
					return nil, errors.New("data too short")
				}
				b.Num = binary.BigEndian.Uint32(cur)
				cur = cur[4:]
			}
		case b.Tag == 4:
			{
				for idx := 0; idx < 8; idx++ {
					if len(cur) < 1 {
						return nil, errors.New("data too short")
					}
					b.Eightbytes[idx] = cur[0]
					cur = cur[1:]
				}
			}
		case b.Tag == 6:
			{
				i := bytes.IndexByte(cur, 0)
				if i < 0 {
					return nil, errors.New("could not parse nul-term string")
				}
				b.String, cur = string(cur[:i]), cur[i+1:]
			}
		}
	}
	return cur, nil
}

func ParseBasic(data []byte) (*Basic, error) {
	b := new(Basic)
	_, err := b.Parse(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}
