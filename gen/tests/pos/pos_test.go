package pos

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
		{Name: "s1_missing_nil", Data: []byte{'a'}},
		{Name: "s2_missing_nil", Data: []byte{'a', 0, 'b'}},
		{Name: "x_short", Data: []byte{'a', 0, 'b', 0, 1}},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := new(Haspos).Parse(c.Data)
			assert.Error(t, err)
		})
	}
}

func TestParseCases(t *testing.T) {
	cases := []struct {
		Name   string
		Data   []byte
		Expect *Haspos
	}{
		{
			Name: "hello_world",
			Data: []byte{
				'H', 'e', 'l', 'l', 'o', 0,
				'W', 'o', 'r', 'l', 'd', '!', 0,
				0, 1, 2, 3,
				'r', 'e', 's', 't',
			},
			Expect: &Haspos{
				S1:   "Hello",
				Pos1: len("Hello") + 1,
				S2:   "World!",
				Pos2: len("Hello World!") + 1,
				X:    0x00010203,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			p := new(Haspos)
			rest, err := p.Parse(c.Data)
			require.NoError(t, err)
			assert.Equal(t, c.Expect, p)
			assert.Equal(t, []byte("rest"), rest)
		})
	}
}
