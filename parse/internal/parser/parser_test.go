package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts := []Option{
		AllowInvalidUTF8(false),
		Debug(false),
		Entrypoint(""),
		MaxExpressions(0),
		Memoize(false),
		Recover(true),
		GlobalStore("foo", "baz"),
		InitState("blah", 42),
		Statistics(&Stats{}, "hmm"),
	}
	src := "const A = 1337;"
	r := strings.NewReader(src)
	_, err := ParseReader("", r, opts...)
	assert.NoError(t, err)
}
