package nulterm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNultermParseLengthErrors(t *testing.T) {
	r := new(Nulterm)
	for n := 0; n < 6; n++ {
		_, err := r.Parse(make([]byte, n))
		require.Error(t, err)
	}
}

func TestNultermMissingNul(t *testing.T) {
	b := []byte{
		1, 2, 3, 4,
		'n', 'o', 'n', 'u', 'l',
	}
	_, err := new(Nulterm).Parse(b)
	assert.Error(t, err)
}

func TestNultermStandard(t *testing.T) {
	n := new(Nulterm)
	b := []byte{
		1, 2, 3, 4,
		'h', 'e', 'l', 'l', 'o', 0,
		5,
		'r', 'e', 's', 't',
	}
	expect := &Nulterm{
		X: 0x01020304,
		S: "hello",
		Y: 5,
	}
	rest, err := n.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, expect, n)
	assert.Equal(t, []byte("rest"), rest)
}
