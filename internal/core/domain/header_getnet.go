package domain

import (
	"strconv"
	"strings"
	"time"
)

type HeaderGetnet struct {
	RegisterType         int8      `txt:"1"`
	ProcessingDate       time.Time `txt:"ddmmyyyy"`
	ProcessingHour       string    `txt:"6"`
	PeriodDate           time.Time `txt:"ddmmyyyy"`
	FileVersion          string    `txt:"8"`
	Headquarter          string    `txt:"15"`
	AcquirerCNPJ         string    `txt:"14"`
	Acquirer             string    `txt:"20"`
	Sequence             int       `txt:"9"`
	AcquirerCode         string    `txt:"2"`
	LayoutVersion        string    `txt:"25"` 
}
func (d HeaderGetnet) GetHeadquarter() int64 {
	n, _ := strconv.ParseInt(strings.TrimSpace(d.Headquarter), 10, 64)
	return n
}
func NewHeaderGetnet() *HeaderGetnet {
	return &HeaderGetnet{}
}
func (d HeaderGetnet) GetProcessingDate() time.Time {
	return d.ProcessingDate
}
func (d HeaderGetnet) GetPeriodInit() time.Time {
	return d.PeriodDate
}
func (d HeaderGetnet) GetPeriodEnd() time.Time {
	return d.PeriodDate
}
func (d HeaderGetnet) GetStatementId() string {
	return "GETNET"
}
func (d HeaderGetnet) GetLayoutVersion() int8 {
	return int8(0)
}
func (d HeaderGetnet) GetAcquirer() string {
	return "GETNET"
}
func (d HeaderGetnet) IsReprocessed() bool {
	return false
}
func (d HeaderGetnet) GetPeriodDates() ([]time.Time, error) {
	ret := make([]time.Time, 1)
	ret[0] = d.PeriodDate
	return ret, nil
}
func (d HeaderGetnet) IsLoaded() bool {
	return d.AcquirerCNPJ == "10440482000154"
}
