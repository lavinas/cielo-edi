package domain

import (
	"fmt"
	"time"
)

type HeaderCielo struct {
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
func NewHeaderCielo() *HeaderCielo {
	return &HeaderCielo{}
}
func (d HeaderCielo) GetHeadquarter() int64 {
	return d.Headquarter
}
func (d HeaderCielo) GetProcessingDate() time.Time {
	return d.ProcessingDate
}
func (d HeaderCielo) GetPeriodInit() time.Time {
	return d.PeriodInit
}
func (d HeaderCielo) GetPeriodEnd() time.Time {
	return d.PeriodEnd
}
func (d HeaderCielo) GetStatementId() string {
	return fmt.Sprintf("%02d", d.StatementId)
}
func (d HeaderCielo) GetLayoutVersion() int8 {
	return d.LayoutVersion
}
func (d HeaderCielo) GetAcquirer() string {
	return d.Acquirer
}
func (d HeaderCielo) IsReprocessed() bool {
	return d.Sequence == 9999999
}
func (d HeaderCielo) GetPeriodDates() ([]time.Time, error) {
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
func (d HeaderCielo) IsLoaded() bool {
	return d.ProcessingDate != time.Time{} && d.Acquirer == "CIELO" && d.LayoutVersion != 0
}
