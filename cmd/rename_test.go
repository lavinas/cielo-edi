package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRename(t *testing.T) {
	err := rename()
	assert.Nil(t, err)
}
