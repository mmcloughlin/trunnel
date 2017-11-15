package eos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEOSParseTooShort(t *testing.T) {
	f := new(Fourbytes)
	for n := 0; n < 4; n++ {
		_, err := f.Parse(make([]byte, n))
		require.Error(t, err)
	}
}

func TestEOSParseTooLong(t *testing.T) {
	f := new(Fourbytes)
	for n := 5; n < 10; n++ {
		_, err := f.Parse(make([]byte, n))
		require.Error(t, err)
	}
}

func TestEOSStandard(t *testing.T) {
	f := new(Fourbytes)
	b := []byte{1, 2, 3, 4}
	rest, err := f.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, &Fourbytes{
		X: 0x0102,
		Y: 0x0304,
	}, f)
	assert.Equal(t, []byte{}, rest)
}
