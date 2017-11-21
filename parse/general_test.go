package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileNotExists(t *testing.T) {
	_, err := File("doesnotexist")
	assert.Error(t, err)
}

func TestErrorReader(t *testing.T) {
	_, err := Reader("", errorReader{})
	assert.Error(t, err)
}

type errorReader struct{}

func (r errorReader) Read(_ []byte) (int, error) {
	return 0, assert.AnError
}
