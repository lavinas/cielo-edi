package string_parser

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

type HeaderCSV struct {
	RegisterType     int8      `txt:"2"`
	HeadquarterId    int64     `txt:"9"`
	ProcDate         time.Time `txt:"ddmmyyyy"`
	PeriodDate       time.Time `txt:"ddmmyyyy"`
	StatementDesc    string    `txt:"39"`
	Acquirer         string    `txt:"8"`
	HeadQquarterName string    `txt:"26"`
	Sequence         int       `txt:"6"`
	ProcessingType   string    `txt:"15"`
	LayoutVersion    string    `txt:"20"`
}


const (
	headerline string = "910238632322021063020210630202106300008358CIELO04I                    014                                                                                                                                                                                 "
	headerlineCsv string = "26,021644942,15052021,14052021,Movimentacao diaria - Cartoes de Debito,Redecard,NESPRESSO PJM             ,000297,DIARIO         ,V1.04 - 07/10 - EEVD"
)

func TestConst(t *testing.T) {
	assert.Equal(t, "txt", tag_name)
}


func TestParseCSVOK(t *testing.T) {
	sp := *NewStringParser("csv")
	header := HeaderCSV{}
	err := sp.Parse(&header, headerlineCsv)
	assert.Nil(t, err)
	assert.Equal(t, int8(26), header.RegisterType)
	assert.Equal(t, int64(21644942), header.HeadquarterId)
	dat, _ := time.Parse("2006-01-02", "2021-05-15")
	assert.Equal(t, dat, header.ProcDate)
	assert.Equal(t, "V1.04 - 07/10 - EEVD", header.LayoutVersion)
}

func TestParseFieldOk(t *testing.T) {
	header := Header{}
	sp := *NewStringParser("position")
	var pos int
	var err error
	// testing int8
	pos, err = sp.ParseField(&header, "RegisterType", headerline, 0)
	assert.Nil(t, err)
	assert.Equal(t, int8(9), header.RegisterType)
	assert.Equal(t, 1, pos)
	// testing int64
	pos, err = sp.ParseField(&header, "HeadquarterId", headerline, 1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1023863232), header.HeadquarterId)
	assert.Equal(t, 11, pos)
	// Testing int
	pos, err = sp.ParseField(&header, "Sequence", headerline, 35)
	assert.Nil(t, err)
	assert.Equal(t, int(8358), header.Sequence)
	assert.Equal(t, 42, pos)
	// Testing String
	pos, err = sp.ParseField(&header, "Acquirer", headerline, 42)
	assert.Nil(t, err)
	assert.Equal(t, "CIELO", header.Acquirer)
	assert.Equal(t, 47, pos)
	// Testing Last Field
	pos, err = sp.ParseField(&header, "Filler", headerline, 73)
	assert.Nil(t, err)
	expVal := "                                                                                                                                                                                 "
	assert.Equal(t, expVal, header.Filler)
	assert.Equal(t, 250, pos)
	// Testing time
	pos, err = sp.ParseField(&header, "ProcDate", headerline, 11)
	assert.Nil(t, err)
	assert.Equal(t, "2021-06-30", header.ProcDate.Format("2006-01-02"))
	assert.Equal(t, 19, pos)
}

func TestParseFieldCSVOk(t *testing.T) {
	header := HeaderCSV{}
	sp := *NewStringParser("csv")
	var pos int
	var err error
	// testing int8
	pos, err = sp.ParseField(&header, "RegisterType", headerlineCsv, 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, pos)
	assert.Equal(t, int8(26), header.RegisterType)
	pos, err = sp.ParseField(&header, "HeadquarterId", headerlineCsv, 1)
	assert.Nil(t, err)
	assert.Equal(t, 2, pos)
	assert.Equal(t, int64(21644942), header.HeadquarterId)
	pos, err = sp.ParseField(&header, "ProcDate", headerlineCsv, 2)
	assert.Nil(t, err)
	assert.Equal(t, 3, pos)
	dat, _ := time.Parse("20060102", "20210515")
	assert.Equal(t, dat, header.ProcDate)
}

func TestParseFieldDataErrors(t *testing.T) {
	header := Header{}
	sp := *NewStringParser("position")
	var pos int
	var err error
	// field name error
	pos, err = sp.ParseField(&header, "RegisterTyp0", headerline, 122)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid field name", err.Error())
	assert.Equal(t, 122, pos)
	//Numeric error
	pos, err = sp.ParseField(&header, "RegisterType", "*10238632322021063020210630202106300008358CIELO04I                    014", 0)
	assert.NotNil(t, err)
	assert.Equal(t, "parsing integer error", err.Error())
	assert.Equal(t, 0, pos)
	//Time error
	pos, err = sp.ParseField(&header, "ProcDate", "910238632322021153020210630202106300008358CIELO04I                    014", 11)
	assert.NotNil(t, err)
	assert.Equal(t, "parsing time \"20211530\": month out of range", err.Error())
	assert.Equal(t, 11, pos)
	//Time error 2
	pos, err = sp.ParseField(&header, "ProcDate", "910238632322021aa3020210630202106300008358CIELO04I                    014", 11)
	assert.NotNil(t, err)
	assert.Equal(t, "parsing time \"2021aa30\" as \"20060102\": cannot parse \"aa30\" as \"01\"", err.Error())
	assert.Equal(t, 11, pos)
	//Unexpected end of line
	pos, err = sp.ParseField(&header, "HeadquarterId", "9102386323", 1)
	assert.NotNil(t, err)
	assert.Equal(t, "unexpected end of txt for parsing this field", err.Error())
	assert.Equal(t, 1, pos)
}

func TestParseAllStruct1(t *testing.T) {
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
	sp := *NewStringParser("position")
	h := Header{}
	pos, err := sp.ParseField(&h, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "tag is not presented", err.Error())
	assert.Equal(t, 0, pos)
}

func TestParseAllStruct2(t *testing.T) {
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
	sp := *NewStringParser("position")
	h := Header{}
	pos, err := sp.ParseField(&h, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "tag is not presented", err.Error())
	assert.Equal(t, 0, pos)
}

func TestInterfaceInvalidStruct(t *testing.T) {
	sp := *NewStringParser("position")
	var errorInterface int
	pos, err := sp.ParseField(&errorInterface, "RegisterType", headerline, 100)
	assert.NotNil(t, err)
	assert.Equal(t, "source interface should be a valid struct", err.Error())
	assert.Equal(t, 100, pos)
	pos, err = sp.ParseField(nil, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "source interface should be a valid struct", err.Error())
	assert.Equal(t, 0, pos)
}

func TestParseIntStruct1(t *testing.T) {
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
	sp := *NewStringParser("position")
	h := Header{}
	pos, err := sp.ParseField(&h, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "not supported field type", err.Error())
	assert.Equal(t, 0, pos)
}

func TestParseIntStruct2(t *testing.T) {
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
	sp := *NewStringParser("position")
	h := Header{}
	pos, err := sp.ParseField(&h, "RegisterType", headerline, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid tag value (should be numeric)", err.Error())
	assert.Equal(t, 0, pos)
}

func TestParseTimeStruct2(t *testing.T) {
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
	sp := *NewStringParser("position")
	h := Header{}
	pos, err := sp.ParseField(&h, "ProcDate", headerline, 11)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid datetime tag value (should be for ex yyyymmdd)", err.Error())
	assert.Equal(t, 11, pos)
}

func TestParseStringStruct3(t *testing.T) {
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
	sp := *NewStringParser("position")
	h := Header{}
	pos, err := sp.ParseField(&h, "ProcDate", "910238632322021-06-3020210630202106300008358CIELO04I                    014", 11)
	assert.Nil(t, err)
	assert.Equal(t, "2021-06-30", h.ProcDate.Format("2006-01-02"))
	assert.Equal(t, 21, pos)
}

func TestParseStringStruct1(t *testing.T) {
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
	sp := *NewStringParser("position")
	h := Header{}
	pos, err := sp.ParseField(&h, "Acquirer", headerline, 42)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid tag value (should be numeric)", err.Error())
	assert.Equal(t, 42, pos)
}

func TestParseStringStruct2(t *testing.T) {
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
	sp := *NewStringParser("position")
	h := Header{}
	pos, err := sp.ParseField(&h, "Acquirer", headerline, 42)
	assert.NotNil(t, err)
	assert.Equal(t, "parsing integer error", err.Error())
	assert.Equal(t, 42, pos)
}

func TestParseOk(t *testing.T) {
	sp := *NewStringParser("position")
	header := Header{}
	err := sp.Parse(&header, headerline)
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

func TestParseOkWithoutFinal(t *testing.T) {
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
	sp := *NewStringParser("position")
	header := Header{}
	err := sp.Parse(&header, headerline)
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

func TestParseErrorShortString(t *testing.T) {
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
	sp := *NewStringParser("position")
	header := Header{}
	err := sp.Parse(&header, "910238632322021063020210630202106300008358CIELO04I                    ")
	assert.NotNil(t, err)
	assert.Equal(t, "LayoutVersion: unexpected end of txt for parsing this field", err.Error())
}

func TestParseErrorInProcDate(t *testing.T) {
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
	sp := *NewStringParser("position")
	header := Header{}
	err := sp.Parse(&header, "910238632322021063020210630202106300008358CIELO04I                    ")
	assert.NotNil(t, err)
	assert.Equal(t, "ProcDate: invalid datetime tag value (should be for ex yyyymmdd)", err.Error())
}

func TestParseOkCsv(t *testing.T) {

}
