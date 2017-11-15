package vararray

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVarArrayParseNoLength(t *testing.T) {
	_, err := new(VarArray).Parse(make([]byte, 1))
	require.Error(t, err)
}

func TestVarArrayParseEmpty(t *testing.T) {
	v := new(VarArray)
	_, err := v.Parse(make([]byte, 2))
	require.NoError(t, err)
	assert.Equal(t, &VarArray{
		NWords: 0,
		Words:  []uint32{},
	}, v)
}

func TestVarArrayTooShort(t *testing.T) {
	v := new(VarArray)
	for n := 1; n < 10; n++ {
		b := make([]byte, 2+4*n-1)
		binary.BigEndian.PutUint16(b, uint16(n))
		_, err := v.Parse(b)
		require.Error(t, err)
	}
}

func TestVarArraySuccess(t *testing.T) {
	v := new(VarArray)
	b := []byte{
		0, 3, // length
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		'r', 'e', 's', 't',
	}
	rest, err := v.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, &VarArray{
		NWords: 3,
		Words: []uint32{
			0x00010203,
			0x04050607,
			0x08090a0b,
		},
	}, v)
	assert.Equal(t, []byte("rest"), rest)
}
