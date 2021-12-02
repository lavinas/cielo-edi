package domain

import (
	"github.com/lavinas/cielo-edi/internal/core/ports"
	"time"
)

type HeaderData struct {
	RegisterType    int8      `parse_txt:"d,1"`
	Headquarter     int64     `parse_txt:"d,10"`
	ProcDate        time.Time `parse_txt:"t,yyyymmdd"`
	PeriodInitDate  time.Time `parse_txt:"t,yyyymmdd"`
	PeriodEndDate   time.Time `parse_txt:"t,yyyymmdd"`
	Sequence        int       `parse_txt:"d,7"`
	Acquirer        string    `parse_txt:"s,5"`
	StatementOption int8      `parse_txt:"d,2"`
	Transmition     string    `parse_txt:"s,1"`
	PostalBox       string    `parse_txt:"s,20"`
	LayoutVersion   int8      `parse_txt:"d, 3"`
	Filler          string    `parse_txt:"s,177"`
}

type Header struct {
	Data   HeaderData
	Parser ports.StringParser
}
