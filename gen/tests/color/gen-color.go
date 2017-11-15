// Code generated by trunnel. DO NOT EDIT.

package color

import "errors"

type Color struct {
	R uint8
	G uint8
	B uint8
}

func (c *Color) Parse(data []byte) ([]byte, error) {
	{
		if len(data) < 1 {
			return nil, errors.New("data too short")
		}
		c.R = data[0]
		data = data[1:]
	}
	{
		if len(data) < 1 {
			return nil, errors.New("data too short")
		}
		c.G = data[0]
		data = data[1:]
	}
	{
		if len(data) < 1 {
			return nil, errors.New("data too short")
		}
		c.B = data[0]
		data = data[1:]
	}
	return data, nil
}