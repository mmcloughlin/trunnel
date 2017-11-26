package fixie

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixieParseTooShort(t *testing.T) {
	f := new(FixieDemo)
	for n := 0; n < 54; n++ {
		_, err := f.Parse(make([]byte, n))
		require.Error(t, err)
	}
}

func TestFixieStandard(t *testing.T) {
	f := new(FixieDemo)
	b := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, // bytes
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', // letters
		0, 1, 2, 3, 4, 5, 6, 7, // shortwords
		0, 1, 2, 3, 4, 5, 6, 7, // words
		0, 1, 2, 3, 4, 5, 6, 7, // big_words[0]
		0, 1, 2, 3, 4, 5, 6, 7, // big_words[1]
		'r', 'g', 'b', // colors[0]
		'R', 'G', 'B', // colors[1]
		'r', 'e', 's', 't',
	}
	rest, err := f.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, &FixieDemo{
		Bytes: [8]uint8{
			0, 1, 2, 3, 4, 5, 6, 7,
		},
		Letters: [8]byte{
			'a', 'b', 'c', 'd',
			'e', 'f', 'g', 'h',
		},
		Shortwords: [4]uint16{
			0x0001, 0x0203,
			0x0405, 0x0607,
		},
		Words: [2]uint32{
			0x00010203,
			0x04050607,
		},
		BigWords: [2]uint64{
			0x0001020304050607,
			0x0001020304050607,
		},
		Colors: [2]*Color{
			{R: 'r', G: 'g', B: 'b'},
			{R: 'R', G: 'G', B: 'B'},
		},
	}, f)
	assert.Equal(t, []byte("rest"), rest)
}
