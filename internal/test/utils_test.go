package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	assert.True(t, FileExists("utils.go"))
	assert.False(t, FileExists("doesnotexist"))
}
