package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type LoggerMock struct {
	lines []string
}

func NewLoggerMock() *LoggerMock {
	return &LoggerMock{}
}

func (l *LoggerMock) Printf(format string, v ...interface{}) {
	l.lines = append(l.lines, fmt.Sprintf(format, v...))
}

func (l *LoggerMock) Println(v ...interface{}) {
	l.lines = append(l.lines, fmt.Sprint(v...))
}

func (l LoggerMock) GetLines() []string {
	return l.lines
}

const (
	path        = "./temp"
	cielofinanc = "010238632322021063020210630202106300008358CIELO04I                    014"
	cielosales  = "010238632322021031020210310202103100008246CIELO03I                    013 "
	cieloant    = "010238632322021051120210511202105110008308CIELO06I                    013"
	redecredit  = "00207022021REDECARDEXTRATO DE MOVIMENTO DE VENDASNESPRESSO PJM         000200021644942DIARIO         V2.01 - 09/06 - EEVC"
	redefin     = "03029092021RedecardExtrato de Movimentacao FinanceiraNESPRESSO PJM         000434021644942DIARIO         V3.01 - 09/06 - EEFI"
	rededebt    = "00,021644942,22092021,21092021,Movimentacao diaria - Cartoes de Debito,Redecard,NESPRESSO PJM             ,000427,DIARIO         ,V1.04 - 07/10 - EEVD"
	getnet      = "02307202106154023072021CEADM1001447355        10440482000154GETNET S.A.         000002146GS                         "
	cieloalelo  = "011013496782020123120201231202012310008177CIELO10I                    013"
)

//------------------------------------------------------------------------------------------

func initPath(path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		endPath(path)
	}
	err := os.Mkdir(path, 0755)
	if err != nil {
		panic(err)
	}
}

func endPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}
}

func createFile(path string, name string, header string) string {
	fn := filepath.Join(path, name)
	os.WriteFile(fn, []byte(header), 0755)
	return fn
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		return false
	}
	return true
}

func TestPeriodsCieloSales(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f1"
	initPath(path)
	createFile(path, "test1.txt", cielosales)
	args := []string{"pm", "periods", "cielovendas", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 1)
	assert.Equal(t, "10/03/2021 - 10/03/2021", result[0])
	endPath(path)
}

func TestGapsCieloSales(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f2"
	initPath(path)
	createFile(path, "test1.txt", cielosales)
	args := []string{"pm", "gaps", "cielovendas", path, "01/03/2021", "30/03/2021"}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "01/03/2021 - 09/03/2021", result[0])
	assert.Equal(t, "11/03/2021 - 30/03/2021", result[1])
	endPath(path)
}

func TestRenameCieloSales(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f3"
	initPath(path)
	fn := createFile(path, "test1.txt", cielosales)
	assert.True(t, fileExists(fn))
	fn = createFile(path, "test2.txt", cielofinanc)
	assert.True(t, fileExists(fn))
	args := []string{"pm", "rename", "cielovendas", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "Yes: test1.txt - CIELO-1023863232-03-2021_03_10-2021_03_10-N-2021_03_10-L013.txt", result[0])
	assert.Equal(t, "No: test2.txt - invalid file", result[1])
	endPath(path)
}

func TestPeriodsRedeSales(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f4"
	initPath(path)
	createFile(path, "test1.txt", redecredit)
	args := []string{"pm", "periods", "redecredito", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 1)
	assert.Equal(t, "07/02/2021 - 07/02/2021", result[0])
	endPath(path)
}

func TestGapsRedeSales(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f5"
	initPath(path)
	createFile(path, "test1.txt", redecredit)
	args := []string{"pm", "gaps", "redecredito", path, "01/01/2021", "31/12/2021"}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "01/01/2021 - 06/02/2021", result[0])
	assert.Equal(t, "08/02/2021 - 31/12/2021", result[1])
	endPath(path)
}

func TestRenameRedeSales(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f6"
	initPath(path)
	fn := createFile(path, "test1.txt", redecredit)
	assert.True(t, fileExists(fn))
	fn = createFile(path, "test2.txt", redefin)
	assert.True(t, fileExists(fn))
	args := []string{"pm", "rename", "redecredito", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "Yes: test1.txt - REDECARD-0021644942-EEVC-2021_02_07-2021_02_07-N-2021_02_07-L002.txt", result[0])
	assert.Equal(t, "No: test2.txt - error parsing", result[1])
	endPath(path)
}

func TestPeriodsRedeFin(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f7"
	initPath(path)
	createFile(path, "test1.txt", redefin)
	args := []string{"pm", "periods", "redefinanceiro", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 1)
	assert.Equal(t, "29/09/2021 - 29/09/2021", result[0])
	endPath(path)
}

func TestGapsRedeFin(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f8"
	initPath(path)
	createFile(path, "test1.txt", redefin)
	args := []string{"pm", "gaps", "redefinanceiro", path, "01/01/2021", "31/12/2021"}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "01/01/2021 - 28/09/2021", result[0])
	assert.Equal(t, "30/09/2021 - 31/12/2021", result[1])
	endPath(path)
}

func TestRenameRedeFin(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f9"
	initPath(path)
	fn := createFile(path, "test1.txt", redefin)
	assert.True(t, fileExists(fn))
	fn = createFile(path, "test2.txt", redecredit)
	assert.True(t, fileExists(fn))
	args := []string{"pm", "rename", "redefinanceiro", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "Yes: test1.txt - REDECARD-0021644942-EEFI-2021_09_29-2021_09_29-N-2021_09_29-L003.txt", result[0])
	assert.Equal(t, "No: test2.txt - error parsing", result[1])
	endPath(path)
}

func TestPeriodsRedeDebt(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f10"
	initPath(path)
	createFile(path, "test1.txt", rededebt)
	args := []string{"pm", "periods", "rededebito", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 1)
	assert.Equal(t, "21/09/2021 - 21/09/2021", result[0])
	endPath(path)
}

func TestGapsRedeDebt(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f11"
	initPath(path)
	createFile(path, "test1.txt", rededebt)
	args := []string{"pm", "gaps", "rededebito", path, "01/01/2021", "31/12/2021"}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "01/01/2021 - 20/09/2021", result[0])
	assert.Equal(t, "22/09/2021 - 31/12/2021", result[1])
	endPath(path)
}

func TestRenameRedeDebt(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f12"
	initPath(path)
	fn := createFile(path, "test1.txt", rededebt)
	assert.True(t, fileExists(fn))
	fn = createFile(path, "test2.txt", redecredit)
	assert.True(t, fileExists(fn))
	args := []string{"pm", "rename", "rededebito", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "Yes: test1.txt - REDECARD-0021644942-EEVD-2021_09_21-2021_09_21-N-2021_09_22-L001.txt", result[0])
	assert.Equal(t, "No: test2.txt - error parsing", result[1])
	endPath(path)
}

func TestPeriodsGetnet(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f13"
	initPath(path)
	createFile(path, "test1.txt", getnet)
	args := []string{"pm", "periods", "getnet", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 1)
	assert.Equal(t, "23/07/2021 - 23/07/2021", result[0])
	endPath(path)
}

func TestGapsGetnet(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f14"
	initPath(path)
	createFile(path, "test1.txt", getnet)
	args := []string{"pm", "gaps", "getnet", path, "01/01/2021", "31/12/2021"}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "01/01/2021 - 22/07/2021", result[0])
	assert.Equal(t, "24/07/2021 - 31/12/2021", result[1])
	endPath(path)
}

func TestRenameGetnet(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f15"
	initPath(path)
	fn := createFile(path, "test1.txt", getnet)
	assert.True(t, fileExists(fn))
	fn = createFile(path, "test2.txt", redecredit)
	assert.True(t, fileExists(fn))
	args := []string{"pm", "rename", "getnet", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "Yes: test1.txt - GETNET-0001447355-GETNET-2021_07_23-2021_07_23-N-2021_07_23-L000.txt", result[0])
	assert.Equal(t, "No: test2.txt - error parsing", result[1])
	endPath(path)
}

func TestPeriodsCieloAlelo(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f1"
	initPath(path)
	createFile(path, "test1.txt", cieloalelo)
	args := []string{"pm", "periods", "cieloalelo", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 1)
	assert.Equal(t, "31/12/2020 - 31/12/2020", result[0])
	endPath(path)
}

func TestGapsCieloAlelo(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f2"
	initPath(path)
	createFile(path, "test1.txt", cieloalelo)
	args := []string{"pm", "gaps", "cieloalelo", path, "01/12/2020", "30/03/2021"}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "01/12/2020 - 30/12/2020", result[0])
	assert.Equal(t, "01/01/2021 - 30/03/2021", result[1])
	endPath(path)
}

func TestRenameCieloAlelo(t *testing.T) {
	logx := NewLoggerMock()
	cm := NewCommandLine(logx)
	path := "./f3"
	initPath(path)
	fn := createFile(path, "test1.txt", cieloalelo)
	assert.True(t, fileExists(fn))
	fn = createFile(path, "test2.txt", cielofinanc)
	assert.True(t, fileExists(fn))
	args := []string{"pm", "rename", "cieloalelo", path}
	err := cm.Run(args)
	assert.Nil(t, err)
	result := logx.GetLines()
	assert.Len(t, result, 2)
	assert.Equal(t, "Yes: test1.txt - CIELO-1101349678-10-2020_12_31-2020_12_31-N-2020_12_31-L013.txt", result[0])
	assert.Equal(t, "No: test2.txt - invalid file", result[1])
	endPath(path)
}
