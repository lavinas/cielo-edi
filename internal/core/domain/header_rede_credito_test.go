package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetStatementId(t *testing.T) {
	header := HeaderRedeCredito{LayoutVersion: "V3.01 - 09/06 - EEFI"}
	id := header.GetStatementId()
	assert.Equal(t, "EEFI", id)
	header = HeaderRedeCredito{LayoutVersion: "V2.01 - 09/06 - EEVC"}
	id = header.GetStatementId()
	assert.Equal(t, "EEVC", id)
}

func TestIsReprocessed(t *testing.T) {
	header := HeaderRedeCredito{ProcessingType: "Diário"}
	rep := header.IsReprocessed()
	assert.False(t, rep)
	header = HeaderRedeCredito{ProcessingType: "Diario"}
	rep = header.IsReprocessed()
	assert.False(t, rep)
	header = HeaderRedeCredito{ProcessingType: "reprocessamento"}
	rep = header.IsReprocessed()
	assert.True(t, rep)
}

func TestGetPeriodDates(t *testing.T) {
	dt, _ := time.Parse("20060102", "20210110")
	header := HeaderRedeCredito{ProcessingDate: dt}
	dates, err := header.GetPeriodDates()
	assert.Nil(t, err)
	assert.Len(t, dates, 1)
	assert.Equal(t, dt, dates[0])
}

