package services

import (
	"errors"
	"io/fs"
	"testing"
	"time"

	"github.com/lavinas/cielo-edi/internal/core/ports"
	"github.com/stretchr/testify/assert"
)

var (
	files = []string{"file1.txt", "file2.txt"}
	path = "./tmp"
)

// Mock of FileInfo
type FileInfoMock struct{
	name string
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
func (f FileManagerMock) GetFiles(string) ([]fs.FileInfo, error){
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
	headquarter int64
	processingDate time.Time
	periodInit time.Time
	periodEnd time.Time
	sequence int
	statementId int8
	layoutVersion int8
	isReprocessed bool
}

func NewHeaderDataMock(hq int64, pd time.Time, pi time.Time, pe time.Time, sq int, st int8, lv int8, ip bool) ports.HeaderDataInterface {
	return &HeaderDataMock{headquarter: hq, processingDate: pd, periodInit: pi, periodEnd: pd, sequence: sq, 
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

// Header mock
type HeaderMock struct {
	headerData ports.HeaderDataInterface
	loaded bool
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
	hd := NewHeaderDataMock(int64(123445), time.Now(), time.Now(), time.Now(), 123, int8(4), int8(14), true)
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
	hd := NewHeaderDataMock(int64(123445), time.Now(), time.Now(), time.Now(), 123, int8(4), int8(14), true)
	he := NewHeaderMock(hd, true)
	// get service
	service := NewService(fm, he)
	err := service.FormatNames(path)
	assert.Nil(t, err)
}