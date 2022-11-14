package controller

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPrefixIdGenerator(t *testing.T) {
	generator := NewPrefixedIdGenerator()
	assert.NotNil(t, generator)

	str, err := generator.Generator("prefix")
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(str, "prefix_"))

	_, err = generator.Generator("")
	assert.Error(t, err)
}
