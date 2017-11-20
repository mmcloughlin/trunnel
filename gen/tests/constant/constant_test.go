package constant

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := map[byte]bool{42: true, 0x42: true, 042: true}
	for x := 0; x < 256; x++ {
		c := new(Constants)
		_, err := c.Parse([]byte{byte(x)})
		if _, ok := cases[byte(x)]; ok {
			require.NoError(t, err)
			assert.Equal(t, &Constants{X: byte(x)}, c)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseEmpty(t *testing.T) {
	_, err := new(Constants).Parse([]byte{})
	assert.Error(t, err)
}
