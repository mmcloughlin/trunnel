package constant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstantValues(t *testing.T) {
	assert.Equal(t, 42, TestDec)
	assert.Equal(t, 0x42, TestHex)
	assert.Equal(t, 042, TestOct)
}
