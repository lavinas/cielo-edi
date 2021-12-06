package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetGapsOk(t *testing.T) {
	init, _ := time.Parse("2006-01-02", "2021-01-01")
	end, _ := time.Parse("2006-01-02", "2021-01-10")
	d := HeaderCielo{PeriodInit: init, PeriodEnd: end}
	dates, err := d.GetPeriodDates()
	assert.Nil(t, err)
	assert.Len(t, dates, 10)
	assert.Contains(t, dates, init)
	assert.Contains(t, dates, end)
	td, _ := time.Parse("2006-01-02", "2021-01-05")
	assert.Contains(t, dates, td)
	td, _ = time.Parse("2006-01-02", "2020-12-31")
	assert.NotContains(t, dates, td)
	td, _ = time.Parse("2006-01-02", "2021-02-01")
	assert.NotContains(t, dates, td)
}

func TestGetGapsEmptyPeriod(t *testing.T) {
	d := HeaderCielo{}
	dates, err := d.GetPeriodDates()
	assert.NotNil(t, err)
	assert.Equal(t, "period is empty", err.Error())
	assert.Equal(t, make([]time.Time, 0), dates)
}

func TestGetGapsWrongPeriod(t *testing.T) {
	init, _ := time.Parse("2006-01-02", "2021-01-10")
	end, _ := time.Parse("2006-01-02", "2021-01-01")
	d := HeaderCielo{PeriodInit: init, PeriodEnd: end}
	dates, err := d.GetPeriodDates()
	assert.NotNil(t, err)
	assert.Equal(t, "initial period after final period", err.Error())
	assert.Equal(t, make([]time.Time, 0), dates)
}

func TestIsValid(t *testing.T) {
	pd, _ := time.Parse("2006-01-02", "2021-01-10")
	d := HeaderCielo{ProcessingDate: pd, Acquirer: "CIELO", LayoutVersion: 14, Statement: "vendas", StatementId: int8(3)}
	iv := d.IsValid()
	assert.True(t, iv)
	d = HeaderCielo{ProcessingDate: pd, Acquirer: "CIELO", LayoutVersion: 14, Statement: "financeiro", StatementId: int8(4)}
	iv = d.IsValid()
	assert.True(t, iv)
	d = HeaderCielo{ProcessingDate: pd, Acquirer: "CIELO", LayoutVersion: 14, Statement: "antecipacoes", StatementId: int8(6)}
	iv = d.IsValid()
	assert.True(t, iv)
	d = HeaderCielo{ProcessingDate: time.Time{}, Acquirer: "CIELO", LayoutVersion: 14, Statement: "vendas", StatementId: int8(3)}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderCielo{ProcessingDate: pd, Acquirer: "REDE", LayoutVersion: 14, Statement: "vendas", StatementId: int8(3)}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderCielo{ProcessingDate: pd, Acquirer: "CIELO", LayoutVersion: 0, Statement: "vendas", StatementId: int8(3)}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderCielo{ProcessingDate: pd, Acquirer: "CIELO", LayoutVersion: 14, Statement: "financeiro", StatementId: int8(3)}
	iv = d.IsValid()
	assert.False(t, iv)
	d = HeaderCielo{ProcessingDate: pd, Acquirer: "CIELO", LayoutVersion: 14, StatementId: int8(3)}
	iv = d.IsValid()
	assert.False(t, iv)

}
