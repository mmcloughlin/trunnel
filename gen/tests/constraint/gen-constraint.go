// Code generated by trunnel. DO NOT EDIT.

package constraint

import (
	"encoding/binary"
	"errors"
)

type Date struct {
	Year  uint16
	Month uint8
	Day   uint8
}

func (d *Date) Parse(data []byte) ([]byte, error) {
	{
		if len(data) < 2 {
			return nil, errors.New("data too short")
		}
		d.Year = binary.BigEndian.Uint16(data)
		if !(1970 <= d.Year && d.Year <= 65535) {
			return nil, errors.New("integer constraint violated")
		}
		data = data[2:]
	}
	{
		if len(data) < 1 {
			return nil, errors.New("data too short")
		}
		d.Month = data[0]
		if !(d.Month == 1 || d.Month == 2 || d.Month == 3 || d.Month == 4 || d.Month == 5 || d.Month == 6 || d.Month == 7 || d.Month == 8 || d.Month == 9 || d.Month == 10 || d.Month == 11 || d.Month == 12) {
			return nil, errors.New("integer constraint violated")
		}
		data = data[1:]
	}
	{
		if len(data) < 1 {
			return nil, errors.New("data too short")
		}
		d.Day = data[0]
		if !(d.Day == 1 || d.Day == 2 || (3 <= d.Day && d.Day <= 31)) {
			return nil, errors.New("integer constraint violated")
		}
		data = data[1:]
	}
	return data, nil
}
