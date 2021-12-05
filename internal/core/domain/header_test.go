package domain

import (
	"testing"

	"github.com/lavinas/cielo-edi/internal/utils/string_parser"
	"github.com/stretchr/testify/assert"
)

var (
	header     Header =  *NewHeader(string_parser.StringParser{})
	headerline string = "910238632322021063020210630202106300008358CIELO04I                    014                                                                                                                                                                                 "
)

func TestParseOkUsingParser(t *testing.T) {
	err := header.Parse(headerline)
	assert.Nil(t, err)
}

func TestParseErrorUsingParser(t *testing.T) {
	err := header.Parse("")
	assert.NotNil(t, err)
	assert.Equal(t, "RegisterType: unexpected end of txt for parsing this field", err.Error())
}
