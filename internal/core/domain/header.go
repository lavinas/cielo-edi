package domain

import (
	"time"

	"github.com/lavinas/cielo-edi/internal/core/ports"
)

type HeaderData struct {
	RegisterType    int8      `parse_txt:"1"`
	Headquarter     int64     `parse_txt:"10"`
	ProcDate        time.Time `parse_txt:"yyyymmdd"`
	PeriodInitDate  time.Time `parse_txt:"yyyymmdd"`
	PeriodEndDate   time.Time `parse_txt:"yyyymmdd"`
	Sequence        int       `parse_txt:"7"`
	Acquirer        string    `parse_txt:"5"`
	StatementOption int8      `parse_txt:"2"`
	Transmition     string    `parse_txt:"1"`
	PostalBox       string    `parse_txt:"20"`
	LayoutVersion   int8      `parse_txt:"3"`
}

type Header struct {
	data   HeaderData
	parser ports.StringParser
}

func NewHeader(parser ports.StringParser) *Header {
	data := HeaderData{} 
	return &Header{data: data, parser: parser}
}

func (h *Header) Parse(txt string) error {
	err := h.parser.Parse(&h.data, txt)
	if err != nil {
		return err
	}
	return nil
}