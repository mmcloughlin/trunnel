package constraint

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func pack(y uint16, m, d uint8) []byte {
	data := []byte{0, 0, m, d}
	binary.BigEndian.PutUint16(data, y)
	return data
}

func TestDateParseLengthErrors(t *testing.T) {
	d := new(Date)
	b := pack(2017, 11, 15)
	for n := 0; n < 4; n++ {
		_, err := d.Parse(b[:n])
		require.Error(t, err)
	}
}

func TestDateParseSuccess(t *testing.T) {
	date := new(Date)
	b := pack(2017, 11, 15)
	extra := []byte("blah")
	rest, err := date.Parse(append(b, extra...))
	require.NoError(t, err)
	assert.Equal(t, &Date{Year: 2017, Month: 11, Day: 15}, date)
	assert.Equal(t, extra, rest)
}

func TestDateParseYearConstraint(t *testing.T) {
	d := new(Date)
	for y := 1700; y < 4000; y++ {
		b := pack(uint16(y), 11, 15)
		_, err := d.Parse(b)
		if y >= 1970 {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, "integer constraint violated")
		}
	}
}

func TestDateParseMonthConstraint(t *testing.T) {
	d := new(Date)
	for m := 0; m < 256; m++ {
		b := pack(2017, byte(m), 15)
		_, err := d.Parse(b)
		if 1 <= m && m <= 12 {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, "integer constraint violated")
		}
	}
}

func TestDateParseDayConstraint(t *testing.T) {
	date := new(Date)
	for d := 0; d < 256; d++ {
		b := pack(2017, 11, byte(d))
		_, err := date.Parse(b)
		if 1 <= d && d <= 31 {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, "integer constraint violated")
		}
	}
}
