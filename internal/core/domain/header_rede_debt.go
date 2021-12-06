package domain

import (
	"strconv"
	"strings"
	"time"
)

var (
	redeDebitoMap = map[string]string{
		"debito": "EEVD",
	}
)

type HeaderRedeDebt struct {
	Statement        string    `txt:"-"`
	RegisterType     int8      `txt:"2"`
	HeadquarterId    int64     `txt:"9"`
	ProcessingDate   time.Time `txt:"ddmmyyyy"`
	PeriodDate       time.Time `txt:"ddmmyyyy"`
	StatementDesc    string    `txt:"39"`
	Acquirer         string    `txt:"8"`
	HeadQquarterName string    `txt:"26"`
	Sequence         int       `txt:"6"`
	ProcessingType   string    `txt:"15"`
	LayoutVersion    string    `txt:"20"`
}

func (d HeaderRedeDebt) GetHeadquarter() int64 {
	return d.HeadquarterId
}
func NewHeaderRedeDebt() *HeaderRedeDebt {
	return &HeaderRedeDebt{}
}
func (d HeaderRedeDebt) GetProcessingDate() time.Time {
	return d.ProcessingDate
}
func (d HeaderRedeDebt) GetPeriodInit() time.Time {
	return d.PeriodDate
}
func (d HeaderRedeDebt) GetPeriodEnd() time.Time {
	return d.PeriodDate
}
func (d HeaderRedeDebt) GetStatementId() string {
	return d.LayoutVersion[16:20]
}
func (d HeaderRedeDebt) GetLayoutVersion() int8 {
	v, err := strconv.ParseInt(string(d.LayoutVersion[1]), 10, 8)
	if err != nil {
		v = int64(0)
	}
	return int8(v)
}
func (d HeaderRedeDebt) GetAcquirer() string {
	return strings.ToUpper(d.Acquirer)
}
func (d HeaderRedeDebt) IsReprocessed() bool {
	return strings.Contains(strings.ToLower(d.ProcessingType), "repro")
}
func (d HeaderRedeDebt) GetPeriodDates() ([]time.Time, error) {
	ret := make([]time.Time, 1)
	ret[0] = d.PeriodDate
	return ret, nil
}
func (d HeaderRedeDebt) IsValid() bool {
	if d.ProcessingDate.Equal(time.Time{}) {
		return false
	}
	if !strings.Contains(strings.ToLower(d.Acquirer), "rede") {
		return false
	}
	if d.LayoutVersion == "" {
		return false
	}
	if _, ok := redeDebitoMap[d.Statement]; !ok {
		return false
	}
	if !strings.Contains(d.LayoutVersion, redeDebitoMap[d.Statement]) {
		return false
	}
	return true
}
