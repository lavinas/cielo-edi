package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Header struct {
	RegisterType    int8      `txt:"1"`        // 0
	HeadquarterId   int64     `txt:"10"`       // 1
	ProcDate        time.Time `txt:"yyyymmdd"` // 11
	PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
	PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
	Sequence        int       `txt:"7"`        // 35
	Acquirer        string    `txt:"5"`        // 42
	StatementOption int8      `txt:"2"`        // 47
	Transmition     string    `txt:"1"`        // 49
	PostalBox       string    `txt:"20"`       // 50
	LayoutVersion   int8      `txt:"3"`        // 70
	Filler          string    `txt:"177"`      // 73
}

const (
	headerline string = "910238632322021063020210630202106300008358CIELO04I                    014                                                                                                                                                                                 "
)

var (
	header Header = Header{}
)

func TestConst(t *testing.T) {
	assert.Equal(t, "txt", tag_name)
}

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
	// Testing time
	pos, err = UnmarshalField(&header, "ProcDate", headerline, 11)
	assert.Nil(t, err)
	assert.Equal(t, "2021-06-30", header.ProcDate.Format("2006-01-02"))
	assert.Equal(t, 19, pos)
}

func TestUnmarshalFieldDataErrors(t *testing.T) {
	var pos int
	var err error
	// field name error
	pos, err = UnmarshalField(&header, "RegisterTyp0", headerline, 122)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid field name", err.Error())
	assert.Equal(t, 122, pos)
	//Numeric error
	pos, err = UnmarshalField(&header, "RegisterType", "*10238632322021063020210630202106300008358CIELO04I                    014", 0)
	assert.NotNil(t, err)
	assert.Equal(t, "parsing integer error", err.Error())
	assert.Equal(t, 0, pos)
	//Time error
	pos, err = UnmarshalField(&header, "ProcDate", "910238632322021153020210630202106300008358CIELO04I                    014", 11)
	assert.NotNil(t, err)
	assert.Equal(t, "parsing time \"20211530\": month out of range", err.Error())
	assert.Equal(t, 11, pos)
	//Time error 2
	pos, err = UnmarshalField(&header, "ProcDate", "910238632322021aa3020210630202106300008358CIELO04I                    014", 11)
	assert.NotNil(t, err)
	assert.Equal(t, "parsing time \"2021aa30\" as \"20060102\": cannot parse \"aa30\" as \"01\"", err.Error())
	assert.Equal(t, 11, pos)
	//Unexpected end of line
	pos, err = UnmarshalField(&header, "HeadquarterId", "9102386323", 1)
	assert.NotNil(t, err)
	assert.Equal(t, "unexpected end of txt for parsing this field", err.Error())
	assert.Equal(t, 1, pos)
}

func TestUnmarshalAllStruct1(t *testing.T) {
	type Header struct {
		RegisterType    int8
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"yyyymmdd"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        string    `txt:"5"`        // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
		Filler          string    `txt:"177"`      // 73
	}
	h := Header{}
	pos, err := UnmarshalField(&h, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "tag is not presented", err.Error())
	assert.Equal(t, 0, pos)
}

func TestUnmarshalAllStruct2(t *testing.T) {
	type Header struct {
		RegisterType    int8      `parsola_txt:"10"`
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"yyyymmdd"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        string    `txt:"5"`        // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
		Filler          string    `txt:"177"`      // 73
	}
	h := Header{}
	pos, err := UnmarshalField(&h, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "tag is not presented", err.Error())
	assert.Equal(t, 0, pos)
}


func TestInterfaceInvalidStruct(t *testing.T) {
	var errorInterface int
	pos, err := UnmarshalField(&errorInterface, "RegisterType", headerline, 100)
	assert.NotNil(t, err)
	assert.Equal(t, "source interface should be a valid struct", err.Error())
	assert.Equal(t, 100, pos)
	pos, err = UnmarshalField(nil, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "source interface should be a valid struct", err.Error())
	assert.Equal(t, 0, pos)
}

func TestUnmarshalIntStruct1(t *testing.T) {
	type Header struct {
		RegisterType    uint16    `txt:"1"`        // 0
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"yyyymmdd"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        string    `txt:"5"`        // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
		Filler          string    `txt:"177"`      // 73
	}
	h := Header{}
	pos, err := UnmarshalField(&h, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "not supported field type", err.Error())
	assert.Equal(t, 0, pos)
}

func TestUnmarshalIntStruct2(t *testing.T) {
	type Header struct {
		RegisterType    int8      `txt:"a"`        // 0
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"yyyymmdd"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        string    `txt:"5"`        // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
		Filler          string    `txt:"177"`      // 73
	}
	h := Header{}
	pos, err := UnmarshalField(&h, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid tag value (should be numeric)", err.Error())
	assert.Equal(t, 0, pos)
}

func TestUnmarshalTimeStruct2(t *testing.T) {
	type Header struct {
		RegisterType    int8      `txt:"1"`        // 0
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"10"`       // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        string    `txt:"5"`        // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
		Filler          string    `txt:"177"`      // 73
	}
	h := Header{}
	pos, err := UnmarshalField(&h, "ProcDate", headerline, 11)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid datetime tag value (should be for ex yyyymmdd)", err.Error())
	assert.Equal(t, 11, pos)
}

func TestUnmarshalStringStruct3(t *testing.T) {
	type Header struct {
		RegisterType    int8      `txt:"1"`          // 0
		HeadquarterId   int64     `txt:"10"`         // 1
		ProcDate        time.Time `txt:"yyyy-mm-dd"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"`   // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"`   // 27
		Sequence        int       `txt:"7"`          // 35
		Acquirer        string    `txt:"ab"`         // 42
		StatementOption int8      `txt:"2"`          // 47
		Transmition     string    `txt:"1"`          // 49
		PostalBox       string    `txt:"20"`         // 50
		LayoutVersion   int8      `txt:"3"`          // 70
		Filler          string    `txt:"177"`        // 73
	}
	h := Header{}
	pos, err := UnmarshalField(&h, "ProcDate", "910238632322021-06-3020210630202106300008358CIELO04I                    014", 11)
	assert.Nil(t, err)
	assert.Equal(t, "2021-06-30", h.ProcDate.Format("2006-01-02"))
	assert.Equal(t, 21, pos)
}

func TestUnmarshalStringStruct1(t *testing.T) {
	type Header struct {
		RegisterType    int8      `txt:"1"`        // 0
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"yyyymmdd"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        string    `txt:"ab"`       // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
		Filler          string    `txt:"177"`      // 73
	}
	h := Header{}
	pos, err := UnmarshalField(&h, "Acquirer", headerline, 42)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid tag value (should be numeric)", err.Error())
	assert.Equal(t, 42, pos)
}

func TestUnmarshalStringStruct2(t *testing.T) {
	type Header struct {
		RegisterType    int8      `txt:"1"`        // 0
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"yyyymmdd"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        int       `txt:"5"`        // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
		Filler          string    `txt:"177"`      // 73
	}
	h := Header{}
	pos, err := UnmarshalField(&h, "Acquirer", headerline, 42)
	assert.NotNil(t, err)
	assert.Equal(t, "parsing integer error", err.Error())
	assert.Equal(t, 42, pos)
}

func TestUnmarshalOk(t *testing.T) {
	header := Header{}
	err := UnMarshal(&header, headerline)
	assert.Nil(t, err)
	assert.Equal(t, int8(9), header.RegisterType)
	assert.Equal(t, int64(1023863232), header.HeadquarterId)
	assert.Equal(t, "20210630", header.ProcDate.Format("20060102"))
	assert.Equal(t, "20210630", header.PeriodInitDate.Format("20060102"))
	assert.Equal(t, "20210630", header.PeriodEndDate.Format("20060102"))
	assert.Equal(t, 8358, header.Sequence)
	assert.Equal(t, "CIELO", header.Acquirer)
	assert.Equal(t, int8(4), header.StatementOption)
	assert.Equal(t, "I", header.Transmition)
	assert.Equal(t, "                    ", header.PostalBox)
	assert.Equal(t, int8(14), header.LayoutVersion)
	blanks := "                                                                                                                                                                                 "
	assert.Equal(t, blanks, header.Filler)
}

func TestUnmarshalOkWithoutFinal(t *testing.T) {
	type Header struct {
		RegisterType    int8      `txt:"1"`        // 0
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"yyyymmdd"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        string    `txt:"5"`        // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
	}
	header := Header{}
	err := UnMarshal(&header, headerline)
	assert.Nil(t, err)
	assert.Equal(t, int8(9), header.RegisterType)
	assert.Equal(t, int64(1023863232), header.HeadquarterId)
	assert.Equal(t, "20210630", header.ProcDate.Format("20060102"))
	assert.Equal(t, "20210630", header.PeriodInitDate.Format("20060102"))
	assert.Equal(t, "20210630", header.PeriodEndDate.Format("20060102"))
	assert.Equal(t, 8358, header.Sequence)
	assert.Equal(t, "CIELO", header.Acquirer)
	assert.Equal(t, int8(4), header.StatementOption)
	assert.Equal(t, "I", header.Transmition)
	assert.Equal(t, "                    ", header.PostalBox)
	assert.Equal(t, int8(14), header.LayoutVersion)
}

func TestUnmarshalErrorShortString(t *testing.T) {
	type Header struct {
		RegisterType    int8      `txt:"1"`        // 0
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"yyyymmdd"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        string    `txt:"5"`        // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
	}
	header := Header{}
	err := UnMarshal(&header, "910238632322021063020210630202106300008358CIELO04I                    ")
	assert.NotNil(t, err)
	assert.Equal(t, "LayoutVersion: unexpected end of txt for parsing this field", err.Error())
}

func TestUnmarshalErrorInProcDate(t *testing.T) {
	type Header struct {
		RegisterType    int8      `txt:"1"`        // 0
		HeadquarterId   int64     `txt:"10"`       // 1
		ProcDate        time.Time `txt:"xxxxyyss"` // 11
		PeriodInitDate  time.Time `txt:"yyyymmdd"` // 19
		PeriodEndDate   time.Time `txt:"yyyymmdd"` // 27
		Sequence        int       `txt:"7"`        // 35
		Acquirer        string    `txt:"5"`        // 42
		StatementOption int8      `txt:"2"`        // 47
		Transmition     string    `txt:"1"`        // 49
		PostalBox       string    `txt:"20"`       // 50
		LayoutVersion   int8      `txt:"3"`        // 70
	}
	header := Header{}
	err := UnMarshal(&header, "910238632322021063020210630202106300008358CIELO04I                    ")
	assert.NotNil(t, err)
	assert.Equal(t, "ProcDate: invalid datetime tag value (should be for ex yyyymmdd)", err.Error())
}