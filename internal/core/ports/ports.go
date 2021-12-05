package ports

import (
	"time"
	"io/fs"
)

type StringParserInterface interface {
	Parse(interface{}, string) error
}

type FileManagerInterface interface {
	GetFiles(string) ([]fs.FileInfo, error)
	GetFirstLine(string, fs.FileInfo) (string, error)
	RenameFile(string, string, string) error
	
}

type HeaderInterface interface {
	Parse(string) error
	IsLoaded() bool
	GetData() HeaderDataInterface	
}

type HeaderDataInterface interface {
	GetHeadquarter() int64
	GetProcessingDate() time.Time
	GetPeriodInit() time.Time
	GetPeriodEnd() time.Time
	GetSequence() int
	GetStatementId() int8
	GetLayoutVersion() int8
	GetAcquirer() string
	IsReprocessed() bool
	GetPeriodDates() ([]time.Time, error) 
}