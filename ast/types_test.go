package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructExtern(t *testing.T) {
	assert.True(t, Struct{}.Extern())
	assert.False(t, Struct{Members: []Member{}}.Extern())
	assert.False(t, Struct{Members: []Member{Field{}}}.Extern())
}
