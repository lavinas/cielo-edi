package services

import (
	"fmt"
	"io/fs"
	"log"

	"github.com/lavinas/cielo-edi/internal/core/ports"
)

const (
	nameFormat string = "CIELO-%010d-%02d-%s-%s-%s-%s-L%03d.txt"
	dateFormat string = "2006_01_02"
)

type Service struct {
	fileManager ports.FileManagerInterface
	header ports.HeaderInterface
}

func NewService(fileManager ports.FileManagerInterface, header ports.HeaderInterface) *Service {
	return &Service{fileManager: fileManager, header: header}
}

func (s *Service) GetHeaderData(path string, file fs.FileInfo) (ports.HeaderDataInterface, error) {
	date, err := s.fileManager.GetFirstLine(path, file)
	if err != nil {
		return nil, err
	}
	if err := s.header.Parse(date); err != nil {
		return nil, fmt.Errorf("error parsing %s", file.Name())
	}
	if !s.header.IsLoaded() {
		return nil, fmt.Errorf("%s is not a valid file", file.Name())
	}
	d := s.header.GetData()
	return d, nil
}

func (s *Service) FormatNames(path string) error {
	log.Printf("Start to rename %s\n", path)

	files, err := s.fileManager.GetFiles(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		header, err := s.GetHeaderData(path, file)
		if err != nil {
			log.Printf("No: %s - %v", file.Name(), err)
			continue
		}
		act := "N"
		if header.IsReprocessed(){
			act = "R"
		}
		newFilename := fmt.Sprintf(nameFormat, header.GetHeadquarter(), header.GetStatementId(), header.GetPeriodInit().Format(dateFormat), 
		header.GetPeriodEnd().Format(dateFormat), act, header.GetProcessingDate().Format(dateFormat), header.GetLayoutVersion())
		err = s.fileManager.RenameFile(path, file.Name(), newFilename)
		if err != nil {
			log.Printf("No: %s - %v", file.Name(), err)
			continue
		}
		log.Printf("Yes: %s", file.Name())
	}
	log.Printf("Finish renaming %s\n", path)
	return nil
}

