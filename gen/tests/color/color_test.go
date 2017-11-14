package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColorParseLengthErrors(t *testing.T) {
	c := new(Color)
	for n := 0; n < 3; n++ {
		_, err := c.Parse(make([]byte, n))
		require.Error(t, err)
	}
}

func TestColorStandard(t *testing.T) {
	c := new(Color)
	b := []byte("Hello World!")
	rest, err := c.Parse(b)
	require.NoError(t, err)
	assert.Equal(t, &Color{R: 'H', G: 'e', B: 'l'}, c)
	assert.Equal(t, []byte("lo World!"), rest)
}
