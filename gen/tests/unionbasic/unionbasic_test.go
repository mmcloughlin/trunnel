package unionbasic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTagConstraint(t *testing.T) {
	b := []byte{42, 0, 0, 0, 0}
	_, err := new(Basic).Parse(b)
	assert.EqualError(t, err, "integer constraint violated")
}

func TestParseEmpty(t *testing.T) {
	_, err := new(Basic).Parse([]byte{})
	assert.Error(t, err)
}

func TestParseShortCases(t *testing.T) {
	tags := []byte{TDate, TInteger, TIntarray, TString}
	for _, tag := range tags {
		_, err := new(Basic).Parse([]byte{byte(tag), 1})
		assert.Error(t, err)
	}
}

func TestParseCases(t *testing.T) {
	cases := []struct {
		Name   string
		Data   []byte
		Expect *Basic
	}{

		{
			Name: "date",
			Data: []byte{2, 7, 225, 11, 15, 'r', 'e', 's', 't'},
			Expect: &Basic{
				Tag: TDate,
				D:   &Date{Year: 2017, Month: 11, Day: 15},
			},
		},
		{
			Name: "integer",
			Data: []byte{3, 0, 1, 2, 3, 'r', 'e', 's', 't'},
			Expect: &Basic{
				Tag: TInteger,
				Num: 0x00010203,
			},
		},
		{
			Name: "int_array",
			Data: []byte{4, 0, 1, 2, 3, 4, 5, 6, 7, 'r', 'e', 's', 't'},
			Expect: &Basic{
				Tag:        TIntarray,
				Eightbytes: [8]byte{0, 1, 2, 3, 4, 5, 6, 7},
			},
		},
		{
			Name: "nulterm",
			Data: []byte{
				6,
				'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd', '!', 0,
				'r', 'e', 's', 't',
			},
			Expect: &Basic{
				Tag:    TString,
				String: "Hello World!",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			u := new(Basic)
			rest, err := u.Parse(c.Data)
			require.NoError(t, err)
			assert.Equal(t, c.Expect, u)
			assert.Equal(t, []byte("rest"), rest)
		})
	}
}

func TestParseDateErrors(t *testing.T) {
	for n := 1; n < 5; n++ {
		b := make([]byte, n)
		b[0] = TDate
		_, err := new(Basic).Parse(b)
		assert.Error(t, err)
	}
}
