package unionlen

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
		{Name: "short_tag", Data: []byte{0}},
		{Name: "short_len", Data: []byte{0, 0, 0}},
		{Name: "short_union", Data: []byte{0, 0, 42, 42}},
		{Name: "short_color_r", Data: []byte{0, 1, 0, 0}},
		{Name: "short_color_g", Data: []byte{0, 1, 0, 1, 'R'}},
		{Name: "short_color_b", Data: []byte{0, 1, 0, 2, 'R', 'G'}},
		{Name: "union_trailing", Data: []byte{0, 1, 0, 6, 'R', 'G', 'B', 'b', 'a', 'd'}},
		{Name: "short_date_year", Data: []byte{0, 2, 0, 1, 0}},
		{Name: "short_date_month", Data: []byte{0, 2, 0, 2, 7, 225}},
		{Name: "short_date_day", Data: []byte{0, 2, 0, 3, 7, 225, 11}},
		{Name: "short_after_union", Data: []byte{0, 1, 0, 3, 'R', 'G', 'B', 0}},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := new(UnionWithLen).Parse(c.Data)
			assert.Error(t, err)
		})
	}
}

func TestParseCases(t *testing.T) {
	cases := []struct {
		Name   string
		Data   []byte
		Expect *UnionWithLen
	}{
		{
			Name: "color",
			Data: []byte{
				0, 1, // color case
				0, 3, // len 3
				'R', 'G', 'B',
				0x13, 0x37,
				'r', 'e', 's', 't',
			},
			Expect: &UnionWithLen{
				Tag:                1,
				UnionLen:           3,
				R:                  'R',
				G:                  'G',
				B:                  'B',
				RightAfterTheUnion: 0x1337,
			},
		},
		{
			Name: "date",
			Data: []byte{
				0, 2, // date case
				0, 7, // should ignore anything over 4 bytes
				7, 225, 11, 16, // the date
				42, 42, 42, // ignored
				0x13, 0x37,
				'r', 'e', 's', 't',
			},
			Expect: &UnionWithLen{
				Tag:                2,
				UnionLen:           7,
				Year:               2017,
				Month:              11,
				Day:                16,
				RightAfterTheUnion: 0x1337,
			},
		},
		{
			Name: "default",
			Data: []byte{
				0x13, 0x37, // should fall to default case
				0, 12,
				'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd', '!',
				0x13, 0x37,
				'r', 'e', 's', 't',
			},
			Expect: &UnionWithLen{
				Tag:                0x1337,
				UnionLen:           12,
				Unparseable:        []byte("Hello World!"),
				RightAfterTheUnion: 0x1337,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			u := new(UnionWithLen)
			rest, err := u.Parse(c.Data)
			require.NoError(t, err)
			assert.Equal(t, c.Expect, u)
			assert.Equal(t, []byte("rest"), rest)
		})
	}
}
