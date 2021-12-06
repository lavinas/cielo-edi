package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRDGetStatementId(t *testing.T) {
	header := HeaderRedeDebt{LayoutVersion: "V3.01 - 09/06 - EEVD"}
	id := header.GetStatementId()
	assert.Equal(t, "EEVD", id)
	header = HeaderRedeDebt{LayoutVersion: "V3.01 - 09/06 - EEVC"}
	id = header.GetStatementId()
	assert.Equal(t, "EEVC", id)
}

func TestRDIsReprocessed(t *testing.T) {
	header := HeaderRedeDebt{ProcessingType: "Diário"}
	rep := header.IsReprocessed()
	assert.False(t, rep)
	header = HeaderRedeDebt{ProcessingType: "Diario"}
	rep = header.IsReprocessed()
	assert.False(t, rep)
	header = HeaderRedeDebt{ProcessingType: "reprocessamento"}
	rep = header.IsReprocessed()
	assert.True(t, rep)
}

func TestRDGetPeriodDates(t *testing.T) {
	dt, _ := time.Parse("20060102", "20210110")
	header := HeaderRedeDebt{PeriodDate: dt}
	dates, err := header.GetPeriodDates()
	assert.Nil(t, err)
	assert.Len(t, dates, 1)
	assert.Equal(t, dt, dates[0])
}

func TestRDIsValidRede(t *testing.T) {
	pd, _ := time.Parse("2006-01-02", "2021-01-10")
	d := HeaderRedeDebt{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "debito"}
	iv := d.IsValid()
	assert.True(t, iv)
	d = HeaderRedeDebt{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVC", Statement: "credito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebt{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVF", Statement: "financeiro"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebt{ProcessingDate: time.Time{}, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "debito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebt{ProcessingDate: pd, Acquirer: "cielo", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "debito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebt{ProcessingDate: pd, Acquirer: "", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "debito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebt{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "", Statement: "debito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebt{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "financeiro"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebt{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: ""}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebt{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVC", Statement: ""}
	iv = d.IsValid()
	assert.False(t, iv)

}
