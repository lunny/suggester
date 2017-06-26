package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeStrings(t *testing.T) {
	var kases = map[string][]string{
		"沪A00561": {"沪", "A", "0", "0", "5", "6", "1"},
	}

	for c, res := range kases {
		assert.Equal(t, res, decodeStrings(string(c)))
	}
}
