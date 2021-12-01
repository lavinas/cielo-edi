package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	var str string = "010238632322021063020210630202106300008358CIELO04I                    014                                                                                                                                                                                 "
	var header Header
	err := header.Parse(str)
	assert.Nil(t, err)
}
