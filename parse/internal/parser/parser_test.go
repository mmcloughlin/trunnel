package parser

import (
	"strings"
	"testing"

	"github.com/mmcloughlin/trunnel/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ParseString(s string) (*ast.File, error) {
	i, err := ParseReader("", strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	return i.(*ast.File), nil
}

func TestConstant(t *testing.T) {
	cases := map[string]string{
		"decimal": `const CONST_ID = 123;`,
		"hex":     `const CONST_ID = 0x7b;`,
		"octal":   `const CONST_ID = 0173;`,
	}
	expect := &ast.File{
		Constants: []*ast.Constant{
			{
				Name:  "CONST_ID",
				Value: 123,
			},
		},
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			f, err := ParseString(src)
			require.NoError(t, err)
			assert.Equal(t, expect, f)
		})
	}
}

func TestExample(t *testing.T) {
	_, err := ParseFile("testdata/example.trunnel")
	require.NoError(t, err)
}
