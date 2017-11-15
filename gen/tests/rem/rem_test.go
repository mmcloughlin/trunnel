package rem

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemParseLengthErrors(t *testing.T) {
	for n := 0; n < 4; n++ {
		_, err := new(Rem).Parse(make([]byte, n))
		require.Error(t, err)
	}
}

// TestRemParseAnyLength confirms that we never get any "rest" bytes returned
// from parsing. These should be taken by the "remaining" array.
func TestRemParseAnyLength(t *testing.T) {
	r := new(Rem)
	for trial := 0; trial < 100; trial++ {
		n := 4 + rand.Intn(1000)
		rest, err := r.Parse(make([]byte, n))
		require.NoError(t, err)
		assert.Equal(t, []byte{}, rest)
		assert.Equal(t, n-4, len(r.Tail))
	}
}

func TestRemParseEmptyTail(t *testing.T) {
	b := []byte{0, 1, 2, 3}
	r := new(Rem)
	rest, err := r.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, rest)
	assert.Equal(t, &Rem{
		Head: 0x00010203,
		Tail: []byte{},
	}, r)
}

func TestRemParseSuccess(t *testing.T) {
	b := []byte{
		0, 1, 2, 3,
		't', 'h', 'e', 't', 'a', 'i', 'l',
	}
	r := new(Rem)
	rest, err := r.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, []byte{}, rest)
	assert.Equal(t, &Rem{
		Head: 0x00010203,
		Tail: []byte("thetail"),
	}, r)
}
