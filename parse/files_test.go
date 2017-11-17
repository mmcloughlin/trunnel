package parse

import (
	"testing"

	"github.com/mmcloughlin/trunnel/test"
	"github.com/stretchr/testify/assert"
)

func TestValidFiles(t *testing.T) {
	test.Glob(t, "./testdata/valid/*.trunnel", valid)
}

func TestFailingFiles(t *testing.T) {
	test.Glob(t, "./testdata/failing/*.trunnel", invalid)
}

func TestTorFiles(t *testing.T) {
	test.Glob(t, "../testdata/tor/*.trunnel", valid)
}

func TestTrunnelFiles(t *testing.T) {
	test.Glob(t, "../testdata/trunnel/*.trunnel", valid)
}

func valid(t *testing.T, filename string) {
	_, err := File(filename)
	assert.NoError(t, err)
}

func invalid(t *testing.T, filename string) {
	_, err := File(filename)
	assert.Error(t, err)
}
