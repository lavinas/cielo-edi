package domain

import (
	"testing"
	"time"


	"github.com/stretchr/testify/assert"
)

func TestRDGetStatementId(t *testing.T) {
	header := HeaderRedeDebito{LayoutVersion: "V3.01 - 09/06 - EEVD"}
	id := header.GetStatementId()
	assert.Equal(t, "EEVD", id)
	header = HeaderRedeDebito{LayoutVersion: "V3.01 - 09/06 - EEVC"}
	id = header.GetStatementId()
	assert.Equal(t, "EEVC", id)
}


func TestRDIsReprocessed(t *testing.T) {
	header := HeaderRedeDebito{ProcessingType: "Diário"}
	rep := header.IsReprocessed()
	assert.False(t, rep)
	header = HeaderRedeDebito{ProcessingType: "Diario"}
	rep = header.IsReprocessed()
	assert.False(t, rep)
	header = HeaderRedeDebito{ProcessingType: "reprocessamento"}
	rep = header.IsReprocessed()
	assert.True(t, rep)
}

func TestRDGetPeriodDates(t *testing.T) {
	dt, _ := time.Parse("20060102", "20210110")
	header := HeaderRedeDebito{PeriodDate: dt}
	dates, err := header.GetPeriodDates()
	assert.Nil(t, err)
	assert.Len(t, dates, 1)
	assert.Equal(t, dt, dates[0])
}


func TestRDIsValidRede(t *testing.T) {
	pd, _ := time.Parse("2006-01-02", "2021-01-10")
	d := HeaderRedeDebito{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "debito"}
	iv := d.IsValid()
	assert.True(t, iv)
	d = HeaderRedeDebito{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVC", Statement: "credito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebito{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVF", Statement: "financeiro"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebito{ProcessingDate: time.Time{}, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "debito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebito{ProcessingDate: pd, Acquirer: "cielo", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "debito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebito{ProcessingDate: pd, Acquirer: "", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "debito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebito{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "", Statement: "debito"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebito{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: "financeiro"}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebito{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVD", Statement: ""}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderRedeDebito{ProcessingDate: pd, Acquirer: "rede", LayoutVersion: "V3.01 - 09/06 - EEVC", Statement: ""}
	iv = d.IsValid()
	assert.False(t, iv)





}