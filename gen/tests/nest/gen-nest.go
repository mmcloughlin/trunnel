// Code generated by trunnel. DO NOT EDIT.

package nest

import "errors"

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

type Rect struct {
	NorthEast *Point
	SouthWest *Point
}

func (r *Rect) Parse(data []byte) ([]byte, error) {
	cur := data
	{
		var err error
		r.NorthEast = new(Point)
		cur, err = r.NorthEast.Parse(cur)
		if err != nil {
			return nil, err
		}
	}
	{
		var err error
		r.SouthWest = new(Point)
		cur, err = r.SouthWest.Parse(cur)
		if err != nil {
			return nil, err
		}
	}
	return cur, nil
}
