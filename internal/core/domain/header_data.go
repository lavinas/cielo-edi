package domain

import (
	"fmt"
	"time"
)

type HeaderData struct {
	RegisterType   int8      `txt:"1"`
	Headquarter    int64     `txt:"10"`
	ProcessingDate time.Time `txt:"yyyymmdd"`
	PeriodInit     time.Time `txt:"yyyymmdd"`
	PeriodEnd      time.Time `txt:"yyyymmdd"`
	Sequence       int       `txt:"7"`
	Acquirer       string    `txt:"5"`
	StatementId    int8      `txt:"2"`
	Transmission   string    `txt:"1"`
	PostalBox      string    `txt:"20"`
	LayoutVersion  int8      `txt:"3"`
}

func (d HeaderData) GetRegisterType() int8 {
	return d.RegisterType
}

func (d HeaderData) GetHeadquarter() int64 {
	return d.Headquarter
}

func (d HeaderData) GetProcessingDate() time.Time {
	return d.ProcessingDate
}

func (d HeaderData) GetPeriodInit() time.Time {
	return d.PeriodInit
}

func (d HeaderData) GetPeriodEnd() time.Time {
	return d.PeriodEnd
}

func (d HeaderData) GetSequence() int {
	return d.Sequence
}

func (d HeaderData) GetAcquirer() string {
	return d.Acquirer
}

func (d HeaderData) GetStatementId() int8 {
	return d.StatementId
}

func (d HeaderData) GetTransmission() string {
	return d.Transmission
}

func (d HeaderData) GetPostalBox() string {
	return d.PostalBox
}

func (d HeaderData) GetLayoutVersion() int8 {
	return d.LayoutVersion
}

func (d HeaderData) IsReprocessed() bool {
	return d.Sequence == 9999999
}

func (d HeaderData) GetPeriodDates() ([]time.Time, error) {
	times := make([]time.Time, 0)
	if d.PeriodInit.Equal(time.Time{}) || d.PeriodEnd.Equal(time.Time{}) {
		return times, fmt.Errorf("period is empty")
	}
	if d.PeriodInit.After(d.PeriodEnd) {
		return times, fmt.Errorf("initial period after final period")
	}
	for t := d.PeriodInit; !t.After(d.PeriodEnd); t = t.Add(24 * time.Hour) {
		times = append(times, t)
	}
	return times, nil
}
