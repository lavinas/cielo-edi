package services

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"
	"time"

	"github.com/lavinas/cielo-edi/internal/core/ports"
	"github.com/stretchr/testify/assert"
)

var (
	files = []string{"file1.txt", "file2.txt"}
	path  = "./tmp"
)

// Mock of FileInfo
type FileInfoMock struct {
	name  string
	isDir bool
}

func NewFileInfoMock(name string, isDir bool) fs.FileInfo {
	return &FileInfoMock{name: name, isDir: isDir}
}
func (i FileInfoMock) Name() string {
	return i.name
}
func (i FileInfoMock) Size() int64 {
	return int64(100)
}
func (i FileInfoMock) Mode() fs.FileMode {
	return fs.FileMode(0)
}
func (i FileInfoMock) ModTime() time.Time {
	return time.Now()
}
func (i FileInfoMock) IsDir() bool {
	return i.isDir
}
func (i FileInfoMock) Sys() interface{} {
	return nil
}

// Mock of filemanager
type FileManagerMock struct {
	files []fs.FileInfo
}

func NewFileManagerMock(files []fs.FileInfo) ports.FileManagerInterface {
	return &FileManagerMock{files: files}
}
func (f FileManagerMock) GetFiles(string) ([]fs.FileInfo, error) {
	return f.files, nil
}
func (f FileManagerMock) GetFirstLine(str string, info fs.FileInfo) (string, error) {
	return str, nil
}
func (f FileManagerMock) RenameFile(string, string, string) error {
	return nil
}

// Header Data Mock
type HeaderDataMock struct {
	headquarter    int64
	processingDate time.Time
	periodInit     time.Time
	periodEnd      time.Time
	sequence       int
	statementId    int8
	layoutVersion  int8
	isReprocessed  bool
}

func NewHeaderDataMock(hq int64, pd time.Time, pi time.Time, pe time.Time, sq int, st int8, lv int8, ip bool) ports.HeaderDataInterface {
	return &HeaderDataMock{headquarter: hq, processingDate: pd, periodInit: pi, periodEnd: pe, sequence: sq,
		statementId: st, layoutVersion: lv, isReprocessed: ip}
}
func (d *HeaderDataMock) GetHeadquarter() int64 {
	return d.headquarter
}
func (d *HeaderDataMock) GetProcessingDate() time.Time {
	return d.processingDate
}
func (d *HeaderDataMock) GetPeriodInit() time.Time {
	return d.periodInit
}
func (d *HeaderDataMock) GetPeriodEnd() time.Time {
	return d.periodEnd
}
func (d *HeaderDataMock) GetSequence() int {
	return d.sequence
}
func (d *HeaderDataMock) GetStatementId() int8 {
	return d.statementId
}
func (d *HeaderDataMock) GetLayoutVersion() int8 {
	return d.layoutVersion
}
func (d *HeaderDataMock) IsReprocessed() bool {
	return d.isReprocessed
}
func (d *HeaderDataMock) GetAcquirer() string {
	return "CIELO"
}
func (d HeaderDataMock) GetPeriodDates() ([]time.Time, error) {
	times := make([]time.Time, 0)
	if d.periodInit.Equal(time.Time{}) || d.periodEnd.Equal(time.Time{}) {
		return times, fmt.Errorf("period is empty")
	}
	if d.periodInit.After(d.periodEnd) {
		return times, fmt.Errorf("initial period after final period")
	}
	for t := d.periodInit; !t.After(d.periodEnd); t = t.Add(24 * time.Hour) {
		times = append(times, t)
	}
	return times, nil
}

// Header mock
type HeaderMock struct {
	headerData ports.HeaderDataInterface
	loaded     bool
}

func NewHeaderMock(hd ports.HeaderDataInterface, lo bool) ports.HeaderInterface {
	return &HeaderMock{headerData: hd, loaded: lo}
}
func (h HeaderMock) Parse(string) error {
	if h.loaded {
		return nil
	}
	return errors.New("Parse Error")
}
func (h HeaderMock) IsLoaded() bool {
	return h.loaded
}
func (h HeaderMock) GetData() ports.HeaderDataInterface {
	return h.headerData
}

func TestGetFilesOk(t *testing.T) {
	// Load FileManager
	fi := make([]fs.FileInfo, 0)
	fi = append(fi, NewFileInfoMock(files[0], false))
	fm := NewFileManagerMock(fi)
	// Load Header
	initDate, _ := time.Parse(printDateFormat, "2021-01-01")
	endDate, _ := time.Parse(printDateFormat, "2021-01-10")
	procDate, _ := time.Parse(printDateFormat, "2021-01-10")
	hd := NewHeaderDataMock(int64(123445), procDate, initDate, endDate, 123, int8(4), int8(14), true)
	he := NewHeaderMock(hd, true)
	// get service
	service := NewService(fm, he)
	rhd, err := service.GetHeaderData("test", fi[0])
	assert.Nil(t, err)
	assert.Equal(t, hd, rhd)
}

func TestFormatNames(t *testing.T) {
	// Load FileManager
	fi := make([]fs.FileInfo, 0)
	fi = append(fi, NewFileInfoMock(files[0], false))
	fm := NewFileManagerMock(fi)
	// Load Header
	initDate, _ := time.Parse(printDateFormat, "2021-01-01")
	endDate, _ := time.Parse(printDateFormat, "2021-01-10")
	procDate, _ := time.Parse(printDateFormat, "2021-01-10")
	hd := NewHeaderDataMock(int64(123445), procDate, initDate, endDate, 123, int8(4), int8(14), true)
	he := NewHeaderMock(hd, true)
	// get service
	service := NewService(fm, he)
	err := service.FormatNames(path)
	assert.Nil(t, err)
}

func TestGetPeriodMap(t *testing.T) {
	// Load FileManager
	fi := make([]fs.FileInfo, 0)
	fi = append(fi, NewFileInfoMock(files[0], false))
	fm := NewFileManagerMock(fi)
	// Load Header
	initDate, _ := time.Parse(printDateFormat, "2021_01_01")
	endDate, _ := time.Parse(printDateFormat, "2021_01_10")
	procDate, _ := time.Parse(printDateFormat, "2021_01_10")
	hd := NewHeaderDataMock(int64(123445), procDate, initDate, endDate, 123, int8(4), int8(14), true)
	he := NewHeaderMock(hd, true)
	// get service
	service := NewService(fm, he)
	dm, err := service.GetPeriodMap(path)
	assert.Nil(t, err)
	assert.Len(t, dm, 10)
	assert.Contains(t, dm, initDate)
	assert.Equal(t, 1, dm[initDate])
	assert.Contains(t, dm, endDate)
	assert.Equal(t, 1, dm[endDate])
	nDate, _ := time.Parse(printDateFormat, "2020_12_31")
	assert.NotContains(t, dm, nDate)
	nDate, _ = time.Parse(printDateFormat, "2020_02_01")
	assert.NotContains(t, dm, nDate)
}

func TestGetPeriodGap(t *testing.T) {
	// Load FileManager
	fi := make([]fs.FileInfo, 0)
	fi = append(fi, NewFileInfoMock(files[0], false))
	fm := NewFileManagerMock(fi)
	// Load Header
	initDate, _ := time.Parse(printDateFormat, "2021_01_01")
	endDate, _ := time.Parse(printDateFormat, "2021_01_10")
	procDate, _ := time.Parse(printDateFormat, "2021_01_10")
	hd := NewHeaderDataMock(int64(123445), procDate, initDate, endDate, 123, int8(4), int8(14), true)
	he := NewHeaderMock(hd, true)
	// get service
	service := NewService(fm, he)
	dates, err := service.GetPeriodGap(path, initDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, dates, 0)
	initDate, _ = time.Parse(printDateFormat, "2020_12_31")
	dates, err = service.GetPeriodGap(path, initDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, dates, 1)
	assert.Contains(t, dates, initDate)
	endDate, _ = time.Parse(printDateFormat, "2021_01_12")
	dates, err = service.GetPeriodGap(path, initDate, endDate)
	assert.Nil(t, err)
	assert.Len(t, dates, 3)
	assert.Contains(t, dates, initDate)
	assert.Contains(t, dates, endDate)
	nDate, _ := time.Parse(printDateFormat, "2020_01_01")
	assert.NotContains(t, dates, nDate)
}

func TestGetPeriod(t *testing.T) {
	// Load FileManager
	fi := make([]fs.FileInfo, 0)
	fi = append(fi, NewFileInfoMock(files[0], false))
	fm := NewFileManagerMock(fi)
	// Load Header
	initDate, _ := time.Parse(printDateFormat, "2021_01_01")
	endDate, _ := time.Parse(printDateFormat, "2021_01_10")
	procDate, _ := time.Parse(printDateFormat, "2021_01_10")
	hd := NewHeaderDataMock(int64(123445), procDate, initDate, endDate, 123, int8(4), int8(14), true)
	he := NewHeaderMock(hd, true)
	// get service
	service := NewService(fm, he)
	dates, err := service.GetPeriod(path)
	assert.Nil(t, err)
	assert.Contains(t, dates, initDate)
	assert.Contains(t, dates, endDate)
}
