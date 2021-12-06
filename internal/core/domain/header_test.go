package domain

import (
	"testing"

	"github.com/lavinas/cielo-edi/internal/utils/string_parser"
	"github.com/stretchr/testify/assert"

)

var (
	cieloheaderline string = "910238632322021063020210630202106300008358CIELO04I                    014                                                                                                                                                                                 "
	redeheaderline string = "03026032021RedecardExtrato de Movimentacao FinanceiraNESPRESSO PJM         000247021644942DIARIO         V3.01 - 09/06 - EEFI"
)

func TestParseCieloOk(t *testing.T) {
	parser := string_parser.NewStringParser("position")
	hd := &HeaderCielo{}
	header := NewHeader(hd, parser)
	err := header.Parse(cieloheaderline)
	assert.Nil(t, err)
}

func TestParseCieloError(t *testing.T) {
	parser := string_parser.NewStringParser("position")
	hd := &HeaderCielo{}
	header := NewHeader(hd, parser)
	err := header.Parse("")
	assert.NotNil(t, err)
	assert.Equal(t, "RegisterType: unexpected end of txt for parsing this field", err.Error())
}

func TestParseRedeOk(t *testing.T) {
	parser := string_parser.NewStringParser("position")
	hd := &HeaderRedeCredito{}
	header := NewHeader(hd, parser)
	err := header.Parse(redeheaderline)
	data := header.GetData()
	assert.Nil(t, err)
	assert.Equal(t, "2021-03-26", data.GetPeriodInit().Format("2006-01-02")) 
}