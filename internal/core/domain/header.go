package domain

import (
	"github.com/lavinas/cielo-edi/internal/core/ports"
)


type Header struct {
	data   ports.HeaderDataInterface
	parser ports.StringParserInterface
}

func NewHeader(parser ports.StringParserInterface) *Header {
	data := &HeaderData{}
	return &Header{data: data, parser: parser}
}

func (h Header) Parse(txt string) error {
	err := h.parser.Parse(h.data, txt)
	if err != nil {
		return err
	}
	return nil
}

func (h Header) GetData() ports.HeaderDataInterface {
	return h.data
}

func (h Header) IsLoaded() bool {
	return h.data.GetAcquirer() == "CIELO"
}

