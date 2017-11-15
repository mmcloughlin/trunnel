// Package gen generates Go parser code from a trunnel AST.
package gen

import (
	"strings"

	"github.com/serenize/snaker"
)

func name(n string) string {
	return snaker.SnakeToCamel(strings.ToLower(n))
}
