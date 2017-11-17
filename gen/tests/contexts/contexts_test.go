package contexts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseErrors(t *testing.T) {
	cases := []struct {
		Name  string
		Flag  Flag
		Count Count
		Data  []byte
	}{
		{Name: "empty", Data: []byte{}},
		{Name: "point_x_constraint", Data: []byte{255}},
		{Name: "point_y_short", Data: []byte{254}},
		{Name: "tsz_0_x_short", Data: []byte{254, 13, 0}},
		{Name: "tsz_0_x_constraint", Data: []byte{254, 13, 0x80, 0, 0, 0}},
		{Name: "tsz_1_y_short", Flag: Flag{1}, Data: []byte{254, 13, 0}},
		{Name: "vsz_a_short", Count: Count{3}, Data: []byte{254, 13, 0, 1, 2, 3}},
		{Name: "vsz_msg_short", Count: Count{2}, Data: []byte{254, 13, 0, 1, 2, 3, 0, 1, 2, 3, 0}},
		{Name: "union_short", Count: Count{2}, Data: []byte{254, 13, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 0}},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := new(Ccomplex).Parse(c.Data, c.Flag, c.Count)
			assert.Error(t, err)
		})
	}
}

func TestParseCases(t *testing.T) {
	cases := []struct {
		Name   string
		Data   []byte
		Flag   Flag
		Count  Count
		Expect *Ccomplex
	}{
		{
			Name:  "flag_0",
			Flag:  Flag{0},
			Count: Count{5},
			Data: []byte{
				42, 13, // point
				0, 1, 2, 3, // x
				4, 5, 6, 7, // a
				'h', 'e', 'l', 'l', 'o', // msg
				'w', 'o', 'r', 'l', 'd', // a
				'r', 'e', 's', 't',
			},
			Expect: &Ccomplex{
				P:   &Point{X: 42, Y: 13},
				Tsz: &Twosize{X: 0x00010203},
				Vsz: &Varsize{A: 0x04050607, Msg: []byte("hello")},
				A:   []byte("world"),
			},
		},
		{
			Name:  "flag_1",
			Flag:  Flag{1},
			Count: Count{6},
			Data: []byte{
				42, 13, // point
				0, 1, // y
				4, 5, 6, 7, // a
				'h', 'e', 'l', 'l', 'o', '!', // msg
				0, 1, 2, 3, 4, 5, // b
				'r', 'e', 's', 't',
			},
			Expect: &Ccomplex{
				P:   &Point{X: 42, Y: 13},
				Tsz: &Twosize{Y: 0x0001},
				Vsz: &Varsize{A: 0x04050607, Msg: []byte("hello!")},
				B:   []uint16{0x0001, 0x0203, 0x0405},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			x := new(Ccomplex)
			rest, err := x.Parse(c.Data, c.Flag, c.Count)
			require.NoError(t, err)
			assert.Equal(t, c.Expect, x)
			assert.Equal(t, []byte("rest"), rest)
		})
	}
}
