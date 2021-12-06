package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedeFinGetStatementId(t *testing.T) {
	header := HeaderRedeFin{LayoutVersion: "V3.01 - 09/06 - EEFI"}
	id := header.GetStatementId()
	assert.Equal(t, "EEFI", id)
	header = HeaderRedeFin{LayoutVersion: "V2.01 - 09/06 - EEVC"}
	id = header.GetStatementId()
	assert.Equal(t, "EEVC", id)
}

func TestRedeFinIsReprocessed(t *testing.T) {
	header := HeaderRedeFin{ProcessingType: "Diário"}
	rep := header.IsReprocessed()
	assert.False(t, rep)
	header = HeaderRedeFin{ProcessingType: "Diario"}
	rep = header.IsReprocessed()
	assert.False(t, rep)
	header = HeaderRedeFin{ProcessingType: "reprocessamento"}
	rep = header.IsReprocessed()
	assert.True(t, rep)
}

func TestRedeFinGetPeriodDates(t *testing.T) {
	dt, _ := time.Parse("20060102", "20210110")
	header := HeaderRedeFin{ProcessingDate: dt}
	dates, err := header.GetPeriodDates()
	assert.Nil(t, err)
	assert.Len(t, dates, 1)
	assert.Equal(t, dt, dates[0])
}

func TestRedeFinIsValidRede(t *testing.T) {
	pd, _ := time.Parse("2006-01-02", "2021-01-10")
	d := HeaderRedeFin{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEFI", Statement: "financeiro"}
	iv := d.IsValid()
	assert.True(t, iv)
	d = HeaderRedeFin{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V2.01 - 09/06 - EEVD", Statement: "debito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeFin{ProcessingDate: pd, Acquirer: "cielo", LayoutVersion: "V3.01 - 09/06 - EEFI", Statement: "financeiro"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeFin{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVC", Statement: "financeiro"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeFin{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "", Statement: "financeiro"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeFin{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "", Statement: ""}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeFin{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "EX", Statement: ""}
	iv = d.IsValid()
	assert.False(t, iv)
}
