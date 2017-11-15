package leftover

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeftoverParseTooShort(t *testing.T) {
	l := new(Leftover)
	for n := 0; n < 16; n++ {
		_, err := l.Parse(make([]byte, n))
		require.Error(t, err)
	}
}

func TestLeftoverParseNonMultiplesOf4(t *testing.T) {
	for n := 1; n < 1000; n++ {
		if n%4 == 0 {
			continue
		}
		_, err := new(Leftover).Parse(make([]byte, n))
		require.Error(t, err)
	}
}

func TestLeftoverParseSuccess(t *testing.T) {
	b := []byte{
		0, 0, 0, 0, // head
		1, 1, 1, 1,
		2, 2, 2, 2, // mid
		3, 3, 3, 3,
		4, 4, 4, 4,
		5, 5, 5, 5, // tail
		6, 6, 6, 6,
	}
	l := new(Leftover)
	rest, err := l.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, rest)
	assert.Equal(t, &Leftover{
		Head: [2]uint32{0x00000000, 0x01010101},
		Mid:  []uint32{0x02020202, 0x03030303, 0x04040404},
		Tail: [2]uint32{0x05050505, 0x06060606},
	}, l)
}
