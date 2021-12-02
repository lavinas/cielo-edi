package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Header struct {
	RegisterType    int8      `parse_txt:"d,1"` // 0
	HeadquarterId   int64     `parse_txt:"d,10"` // 1
	ProcDate        time.Time `parse_txt:"t,yyyymmdd"` // 11
	PeriodInitDate  time.Time `parse_txt:"t,yyyymmdd"` // 19
	PeriodEndDate   time.Time `parse_txt:"t,yyyymmdd"` // 27
	Sequence        int       `parse_txt:"d,7"`        // 35
	Acquirer        string    `parse_txt:"s,5"`        // 42
	StatementOption int8      `parse_txt:"d,2"`        // 47
	Transmition     string    `parse_txt:"s,1"`        // 49
	PostalBox       string    `parse_txt:"s,20"`       // 50
	LayoutVersion   int8      `parse_txt:"d, 3"`       // 70
	Filler          string    `parse_txt:"s,177"`      // 73
}

const (
	headerline string = "910238632322021063020210630202106300008358CIELO04I                    014                                                                                                                                                                                 "
)

var (
	header Header = Header{}
)

func TestUnmarshalFieldOk(t *testing.T) {
	var pos int
	var err error
	// testing int8
	pos, err = UnmarshalField(&header, "RegisterType", headerline, 0)
	assert.Nil(t, err)
	assert.Equal(t, int8(9), header.RegisterType)
	assert.Equal(t, 1, pos)
	// testing int64
	pos, err = UnmarshalField(&header, "HeadquarterId", headerline, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1023863232), header.HeadquarterId)
	assert.Equal(t, 11, pos)
	// Testing int
	pos, err = UnmarshalField(&header, "Sequence", headerline, 35)
	assert.Nil(t, err)
	assert.Equal(t, int(8358), header.Sequence)
	assert.Equal(t, 42, pos)
	// Testing String
	pos, err = UnmarshalField(&header, "Acquirer", headerline, 42)
	assert.Nil(t, err)
	assert.Equal(t, "CIELO", header.Acquirer)
	assert.Equal(t, 47, pos)
	// Testing Last Field
	pos, err = UnmarshalField(&header, "Filler", headerline, 73)
	assert.Nil(t, err)
	expVal := "                                                                                                                                                                                 "
	assert.Equal(t, expVal, header.Filler)
	assert.Equal(t, 250, pos)
	// Testing field error
	pos, err = UnmarshalField(&header, "Fill", headerline, 73)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid field name", err.Error())
	assert.Equal(t, 73, pos)
	// Testing field error
	pos, err = UnmarshalField(&header, "Fill", headerline, 73)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid field name", err.Error())
	assert.Equal(t, 73, pos)
}