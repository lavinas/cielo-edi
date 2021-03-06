package ports

import (
	"io/fs"
	"time"
)

type StringParserInterface interface {
	Parse(interface{}, string) error
}

type FileManagerInterface interface {
	GetFiles(string) ([]fs.FileInfo, error)
	GetFirstLine(string, fs.FileInfo) (string, error)
	RenameFile(string, string, string) error
}

type LoggerInterface interface {
	Printf(string, ...interface{})
	Println(...interface{})
}

type HeaderInterface interface {
	Parse(string) error
	IsValid() bool
	GetData() HeaderDataInterface
}

type HeaderDataInterface interface {
	GetHeadquarter() int64
	GetProcessingDate() time.Time
	GetPeriodInit() time.Time
	GetPeriodEnd() time.Time
	GetStatementId() string
	GetLayoutVersion() int8
	GetAcquirer() string
	IsReprocessed() bool
	GetPeriodDates() ([]time.Time, error)
	IsValid() bool
}

type ServiceInterface interface {
	FormatNames(string) ([]string, error)
	GetGapGrouped(string, time.Time, time.Time) ([]string, error)
	GetPeriodGrouped(string) ([]string, error)
}

type CommandLineInterface interface {
	Run([]string) error
}
