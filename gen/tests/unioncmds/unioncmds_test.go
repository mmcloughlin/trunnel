package unioncmds

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
		{Name: "fail_case", Data: []byte{2}},
		{Name: "default_short", Data: []byte{123, 0, 1, 2}},
		{Name: "y_short", Data: []byte{1, 0, 1, 2}},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := new(UnionCmds).Parse(c.Data)
			assert.Error(t, err)
		})
	}
}

func TestParseCases(t *testing.T) {
	cases := []struct {
		Name   string
		Data   []byte
		Expect *UnionCmds
	}{
		{
			Name: "ignore",
			Data: []byte{
				1, // ignore case
				0, 1, 2, 3,
				'r', 'e', 's', 't',
			},
			Expect: &UnionCmds{
				Tag: 1,
				Y:   0x00010203,
			},
		},
		{
			Name: "default",
			Data: []byte{
				42,
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				'r', 'e', 's', 't',
			},
			Expect: &UnionCmds{
				Tag: 42,
				X: [2]uint32{
					0x00010203,
					0x04050607,
				},
				Y: 0x08090a0b,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			u := new(UnionCmds)
			rest, err := u.Parse(c.Data)
			require.NoError(t, err)
			assert.Equal(t, c.Expect, u)
			assert.Equal(t, []byte("rest"), rest)
		})
	}
}
