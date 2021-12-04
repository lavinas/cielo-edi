package domain

import (
	"encoding/json"
	"time"

	"github.com/lavinas/cielo-edi/internal/core/ports"
)

type HeaderData struct {
	RegisterType   int8      `txt:"1" json:"register_type"`
	Headquarter    int64     `txt:"10" json:"headquarter"`
	ProcessingDate time.Time `txt:"yyyymmdd" json:"processing_date"`
	PeriodInit     time.Time `txt:"yyyymmdd" json:"period_init"`
	PeriodEnd      time.Time `txt:"yyyymmdd" json:"period_end"`
	Sequence       int       `txt:"7" json:"sequence"`
	Acquirer       string    `txt:"5" json:"acquirer"`
	StatementId    int8      `txt:"2" json:"statement_id"`
	Transmition    string    `txt:"1"  json:"transmission"`
	PostalBox      string    `txt:"20"  json:"postal_box"`
	LayoutVersion  int8      `txt:"3" json:"layout_version"`
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

func (h *Header) GetJsonData() ([]byte, error) {
	json, err := json.Marshal(h.data)
	if err != nil {
		return make([]byte, 0), err
	}
	return json, nil
}

func (h *Header) GetData() HeaderData {
	return h.data
}

func (h *Header) IsLoaded() bool {
	return h.data.Acquirer == "CIELO"
}
