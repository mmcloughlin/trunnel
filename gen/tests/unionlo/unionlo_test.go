package unionlo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseErrors(t *testing.T) {
	cases := []struct {
		Name string
		Data []byte
	}{
		{Name: "empty", Data: []byte{}},
		{Name: "short_union", Data: []byte{1}},
		{Name: "short_1", Data: []byte{1, 0, 0, 0, 0, 0, 0, 0, 0}},
		{Name: "short_4_byte", Data: []byte{4, 0, 0, 0, 0, 0, 0, 0, 0}},
		{Name: "short_4_z", Data: []byte{4, 42, 1, 0, 0, 0, 0, 0, 0, 0, 0}},
		{Name: "leftoverlen", Data: []byte{4, 42, 1, 2, 255, 0, 0, 0, 0, 0, 0, 0}},
		{Name: "trailing_data", Data: []byte{1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0}},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := new(Unlo).Parse(c.Data)
			assert.Error(t, err)
		})
	}
}

func TestParseCases(t *testing.T) {
	cases := []struct {
		Name   string
		Data   []byte
		Expect *Unlo
	}{
		{
			Name: "x",
			Data: []byte{
				1,
				42,
				3,
				'a', 'b', 'c',
				'r', 'e', 's', 't',
			},
			Expect: &Unlo{
				Tag:         1,
				X:           42,
				Leftoverlen: 3,
				Leftovers:   []byte("abc"),
			},
		},
		{
			Name: "y",
			Data: []byte{
				2,
				'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd', '!',
				3,
				'a', 'b', 'c',
				'r', 'e', 's', 't',
			},
			Expect: &Unlo{
				Tag:         2,
				Y:           []byte("Hello World!"),
				Leftoverlen: 3,
				Leftovers:   []byte("abc"),
			},
		},
		{
			Name: "z",
			Data: []byte{
				4,
				42, 0, 1, 2, 3, 4, 5,
				3,
				'a', 'b', 'c',
				'r', 'e', 's', 't',
			},
			Expect: &Unlo{
				Tag:         4,
				Byte:        42,
				Z:           []uint16{0x0001, 0x0203, 0x0405},
				Leftoverlen: 3,
				Leftovers:   []byte("abc"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			u := new(Unlo)
			rest, err := u.Parse(c.Data)
			require.NoError(t, err)
			assert.Equal(t, c.Expect, u)
			assert.Equal(t, []byte("rest"), rest)
		})
	}
}
