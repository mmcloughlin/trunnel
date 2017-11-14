package ints

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntsParseLengthErrors(t *testing.T) {
	x := new(Ints)
	for n := 0; n < 15; n++ {
		_, err := x.Parse(make([]byte, n))
		require.Error(t, err)
	}
}

func TestIntsStandard(t *testing.T) {
	x := new(Ints)
	b := []byte{
		1,
		2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		'r', 'e', 's', 't',
	}
	expect := &Ints{
		Byte:  0x01,
		Word:  0x0203,
		Dword: 0x04050607,
		Qword: 0x08090a0b0c0d0e0f,
	}
	rest, err := x.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, expect, x)
	assert.Equal(t, []byte("rest"), rest)
}
