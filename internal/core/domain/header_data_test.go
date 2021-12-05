package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func TestGetGapsOk(t *testing.T) {
	init, _ := time.Parse("2006-01-02", "2021-01-01")
	end,  _ := time.Parse("2006-01-02", "2021-01-10")
	d := HeaderData{PeriodInit: init, PeriodEnd: end}
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
	d := HeaderData{}
	dates, err := d.GetPeriodDates()
	assert.NotNil(t, err)
	assert.Equal(t, "period is empty", err.Error())
	assert.Equal(t, make([]time.Time, 0), dates)
}

func TestGetGapsWrongPeriod(t *testing.T) {
	init, _ := time.Parse("2006-01-02", "2021-01-10")
	end,  _ := time.Parse("2006-01-02", "2021-01-01")
	d := HeaderData{PeriodInit: init, PeriodEnd: end}
	dates, err := d.GetPeriodDates()
	assert.NotNil(t, err)
	assert.Equal(t, "initial period after final period", err.Error())
	assert.Equal(t, make([]time.Time, 0), dates)
}
