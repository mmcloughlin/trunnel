package nest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNestParseLengthErrors(t *testing.T) {
	r := new(Rect)
	for n := 0; n < 4; n++ {
		_, err := r.Parse(make([]byte, n))
		require.Error(t, err)
	}
}

func TestNestStandard(t *testing.T) {
	r := new(Rect)
	b := []byte{
		1, 2,
		3, 4,
		'r', 'e', 's', 't',
	}
	expect := &Rect{
		NorthEast: &Point{X: 1, Y: 2},
		SouthWest: &Point{X: 3, Y: 4},
	}
	rest, err := r.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, expect, r)
	assert.Equal(t, []byte("rest"), rest)
}
