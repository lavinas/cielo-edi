package domain

import (
	"strconv"
	"strings"
	"time"
)

var (
	redeFinMap = map[string]string{
		"financeiro": "EEFI",
	}
)

type HeaderRedeFin struct {
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

func (d HeaderRedeFin) GetHeadquarter() int64 {
	return d.Headquarter
}
func NewHeaderRedeFin() *HeaderRedeFin {
	return &HeaderRedeFin{}
}
func (d HeaderRedeFin) GetProcessingDate() time.Time {
	return d.ProcessingDate
}
func (d HeaderRedeFin) GetPeriodInit() time.Time {
	return d.ProcessingDate
}
func (d HeaderRedeFin) GetPeriodEnd() time.Time {
	return d.ProcessingDate
}
func (d HeaderRedeFin) GetStatementId() string {
	return d.LayoutVersion[16:20]
}
func (d HeaderRedeFin) GetLayoutVersion() int8 {
	v, err := strconv.ParseInt(string(d.LayoutVersion[1]), 10, 8)
	if err != nil {
		v = int64(0)
	}
	return int8(v)
}
func (d HeaderRedeFin) GetAcquirer() string {
	return strings.ToUpper(d.Acquirer)
}
func (d HeaderRedeFin) IsReprocessed() bool {
	return strings.Contains(strings.ToLower(d.ProcessingType), "repro")
}
func (d HeaderRedeFin) GetPeriodDates() ([]time.Time, error) {
	ret := make([]time.Time, 1)
	ret[0] = d.ProcessingDate
	return ret, nil
}

func (d HeaderRedeFin) IsValid() bool {
	if d.ProcessingDate.Equal(time.Time{}) {
		return false
	}
	if !strings.Contains(strings.ToLower(d.Acquirer), "rede") {
		return false
	}
	if d.LayoutVersion == "" {
		return false
	}
	if _, ok := redeFinMap[d.Statement]; !ok {
		return false
	}
	if !strings.Contains(d.LayoutVersion, redeFinMap[d.Statement]) {
		return false
	}
	return true
}
