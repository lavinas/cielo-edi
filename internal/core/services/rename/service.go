package rename

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/lavinas/cielo-edi/internal/core/domain"
)

const (
	nameFormat string = "CIELO-%010d-%02d-%s-%s-%s-%s-L%03d.txt"
	dateFormat string = "2006_01_02"
)

type RenameService struct {
	header domain.Header
}

func (s *RenameService) FormatFilesName(path string) error {
	log.Printf("Start to rename %s\n", path)

	files, err := getFiles(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			log.Printf("No: dir %s\n", file.Name())
			continue
		}
		buf, err := readFile(path, file.Name())
		if err != nil {
			log.Printf("No: read error %s\n", file.Name())
			continue
		}
		if !buf.Scan() {
			log.Printf("No: scan error %s\n", file.Name())
			continue
		}
		hDate := buf.Text()
		if err := s.header.Parse(hDate); err != nil {
			log.Printf("No: parse error %s", file.Name())
			continue
		}
		if !s.header.IsLoaded() {
			log.Printf("No: no cielo %s", file.Name())
			continue
		}
		if err := renameFile(path, file.Name(), s.header); err != nil {
			log.Printf("No: rename error %s", file.Name())
			continue
		}
		log.Printf("Yes: %s", file.Name())
	}
	log.Printf("Finish renaming %s\n", path)
	return nil
}

func NewRenameService(header domain.Header) *RenameService {
	return &RenameService{header: header}
}

func getFiles(path string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return make([]fs.FileInfo, 0), err
	}
	return files, nil
}

func readFile(path string, filename string) (*bufio.Scanner, error) {
	fileIO, err := os.OpenFile(filepath.Join(path, filename), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(fileIO)
	return scanner, nil
}

func renameFile(path string, filename string, header domain.Header) error {
	d := header.GetData()
	init := d.PeriodInit.Format(dateFormat)
	end := d.PeriodEnd.Format(dateFormat)
	proc := d.ProcessingDate.Format(dateFormat)
	act := "N"
	if d.Sequence == 9999999 {
		act = "R"
	}
	newFilename := fmt.Sprintf(nameFormat, d.Headquarter, d.StatementId, init, end, act, proc, d.LayoutVersion)
	from := filepath.Join(path, filename)
	to := filepath.Join(path, newFilename)
	if err := os.Rename(from, to); err != nil {
		return err
	}
	return nil
}
