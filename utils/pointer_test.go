package utils_test

import (
	"testing"

	"github.com/dmitrymomot/solana/utils"
	"github.com/stretchr/testify/assert"
)

func TestPointer(t *testing.T) {
	var i int = 123
	assert.IsType(t, &i, utils.Pointer(i))
	assert.Equal(t, &i, utils.Pointer(i))

	var s string = "abc"
	assert.IsType(t, &s, utils.Pointer(s))
	assert.Equal(t, &s, utils.Pointer(s))

	var b bool = true
	assert.IsType(t, &b, utils.Pointer(b))
	assert.Equal(t, &b, utils.Pointer(b))
}
