package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunUrl(t *testing.T) {

	result := runUrl("http://127.0.0.1", 1)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.successCount)
}
