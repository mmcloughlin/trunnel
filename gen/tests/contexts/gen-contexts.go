// Code generated by trunnel. DO NOT EDIT.

package contexts

import (
	"encoding/binary"
	"errors"
)

type Flag struct {
	Flagval uint8
}

type Count struct {
	Countval uint8
}

type Point struct {
	X uint8
	Y uint8
}

func (p *Point) Parse(data []byte) ([]byte, error) {
	cur := data
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		p.X = cur[0]
		if !(0 <= p.X && p.X <= 254) {
			return nil, errors.New("integer constraint violated")
		}
		cur = cur[1:]
	}
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		p.Y = cur[0]
		cur = cur[1:]
	}
	return cur, nil
}

type Twosize struct {
	X uint32
	Y uint16
}

func (t *Twosize) Parse(data []byte, flag Flag) ([]byte, error) {
	cur := data
	{
		switch {
		case flag.Flagval == 0:
			{
				if len(cur) < 4 {
					return nil, errors.New("data too short")
				}
				t.X = binary.BigEndian.Uint32(cur)
				if !(0 <= t.X && t.X <= 2147483647) {
					return nil, errors.New("integer constraint violated")
				}
				cur = cur[4:]
			}
		case flag.Flagval == 1 || flag.Flagval == 3:
			{
				if len(cur) < 2 {
					return nil, errors.New("data too short")
				}
				t.Y = binary.BigEndian.Uint16(cur)
				cur = cur[2:]
			}
		}
	}
	return cur, nil
}

type Varsize struct {
	A   uint32
	Msg []uint8
}

func (v *Varsize) Parse(data []byte, count Count) ([]byte, error) {
	cur := data
	{
		if len(cur) < 4 {
			return nil, errors.New("data too short")
		}
		v.A = binary.BigEndian.Uint32(cur)
		cur = cur[4:]
	}
	{
		v.Msg = make([]uint8, int(count.Countval))
		for idx := 0; idx < int(count.Countval); idx++ {
			if len(cur) < 1 {
				return nil, errors.New("data too short")
			}
			v.Msg[idx] = cur[0]
			cur = cur[1:]
		}
	}
	return cur, nil
}

type Ccomplex struct {
	P   *Point
	Tsz *Twosize
	Vsz *Varsize
	A   []uint8
	B   []uint16
}

func (c *Ccomplex) Parse(data []byte, flag Flag, count Count) ([]byte, error) {
	cur := data
	{
		var err error
		c.P = new(Point)
		cur, err = c.P.Parse(cur)
		if err != nil {
			return nil, err
		}
	}
	{
		var err error
		c.Tsz = new(Twosize)
		cur, err = c.Tsz.Parse(cur, flag)
		if err != nil {
			return nil, err
		}
	}
	{
		var err error
		c.Vsz = new(Varsize)
		cur, err = c.Vsz.Parse(cur, count)
		if err != nil {
			return nil, err
		}
	}
	{
		if len(cur) < int(count.Countval) {
			return nil, errors.New("data too short")
		}
		restore := cur[int(count.Countval):]
		cur = cur[:int(count.Countval)]
		switch {
		case flag.Flagval == 0:
			{
				c.A = make([]uint8, 0)
				for len(cur) > 0 {
					var tmp uint8
					if len(cur) < 1 {
						return nil, errors.New("data too short")
					}
					tmp = cur[0]
					cur = cur[1:]
					c.A = append(c.A, tmp)
				}
			}
		case flag.Flagval == 1:
			{
				c.B = make([]uint16, 0)
				for len(cur) > 0 {
					var tmp uint16
					if len(cur) < 2 {
						return nil, errors.New("data too short")
					}
					tmp = binary.BigEndian.Uint16(cur)
					cur = cur[2:]
					c.B = append(c.B, tmp)
				}
			}
		}
		if len(cur) > 0 {
			return nil, errors.New("trailing data disallowed")
		}
		cur = restore
	}
	return cur, nil
}