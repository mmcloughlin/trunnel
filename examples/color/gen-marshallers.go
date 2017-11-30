// Code generated by trunnel. DO NOT EDIT.

package color

import "errors"

type Color struct {
	R uint8
	G uint8
	B uint8
}

func (c *Color) Parse(data []byte) ([]byte, error) {
	cur := data
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		c.R = cur[0]
		cur = cur[1:]
	}
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		c.G = cur[0]
		cur = cur[1:]
	}
	{
		if len(cur) < 1 {
			return nil, errors.New("data too short")
		}
		c.B = cur[0]
		cur = cur[1:]
	}
	return cur, nil
}

func ParseColor(data []byte) (*Color, error) {
	c := new(Color)
	_, err := c.Parse(data)
	if err != nil {
		return nil, err
	}
	return c, nil
}