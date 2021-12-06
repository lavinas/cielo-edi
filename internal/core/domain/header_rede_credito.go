package domain

import (
	"strconv"
	"strings"
	"time"
)

var (
	redeCreditoMap = map[string]string{
		"credito":    "EEVC",
		"financeiro": "EEFI",
	}
)

type HeaderRedeCredito struct {
	Statement            string    `txt:"-"`
	RegisterType         int8      `txt:"3"`
	ProcessingDate       time.Time `txt:"ddmmyyyy"`
	Acquirer             string    `txt:"8"`
	StatementDescription string    `txt:"34"`
	HeadquarterName      string    `txt:"22"`
	Sequence             int       `txt:"6"`
	Headquarter          int64     `txt:"9"`
	ProcessingType       string    `txt:"15"`
	LayoutVersion        string    `txt:"20"`
}

func (d HeaderRedeCredito) GetHeadquarter() int64 {
	return d.Headquarter
}
func NewHeaderRedeCredito() *HeaderRedeCredito {
	return &HeaderRedeCredito{}
}
func (d HeaderRedeCredito) GetProcessingDate() time.Time {
	return d.ProcessingDate
}
func (d HeaderRedeCredito) GetPeriodInit() time.Time {
	return d.ProcessingDate
}
func (d HeaderRedeCredito) GetPeriodEnd() time.Time {
	return d.ProcessingDate
}
func (d HeaderRedeCredito) GetStatementId() string {
	return d.LayoutVersion[16:20]
}
func (d HeaderRedeCredito) GetLayoutVersion() int8 {
	v, err := strconv.ParseInt(string(d.LayoutVersion[1]), 10, 8)
	if err != nil {
		v = int64(0)
	}
	return int8(v)
}
func (d HeaderRedeCredito) GetAcquirer() string {
	return d.Acquirer
}
func (d HeaderRedeCredito) IsReprocessed() bool {
	return strings.Contains(strings.ToLower(d.ProcessingType), "repro")
}
func (d HeaderRedeCredito) GetPeriodDates() ([]time.Time, error) {
	ret := make([]time.Time, 1)
	ret[0] = d.ProcessingDate
	return ret, nil
}

func (d HeaderRedeCredito) IsValid() bool {
	if d.ProcessingDate.Equal(time.Time{}) {
		return false
	}
	if !strings.Contains(strings.ToLower(d.Acquirer), "rede") {
		return false
	}
	if d.LayoutVersion == "" {
		return false
	}
	if _, ok := redeCreditoMap[d.Statement]; !ok {
		return false
	}
	if !strings.Contains(d.LayoutVersion, redeCreditoMap[d.Statement]) {
		return false
	}
	return true
}
