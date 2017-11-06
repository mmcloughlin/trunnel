package trunnel

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestAdhoc(t *testing.T) {
	i, err := ParseFile("example.trunnel")
	spew.Dump(i)
	spew.Dump(err)
}
