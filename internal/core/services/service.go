package services

import (
	"fmt"
	"io/fs"
	"log"
	"sort"
	"time"

	"github.com/lavinas/cielo-edi/internal/core/ports"
)

const (
	nameFormat      string = "CIELO-%010d-%02d-%s-%s-%s-%s-L%03d.txt"
	printDateFormat string = "2006_01_02"
)

type Service struct {
	fileManager ports.FileManagerInterface
	header      ports.HeaderInterface
}

func NewService(fileManager ports.FileManagerInterface, header ports.HeaderInterface) *Service {
	return &Service{fileManager: fileManager, header: header}
}

func (s Service) GetHeaderData(path string, file fs.FileInfo) (ports.HeaderDataInterface, error) {
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

func (s Service) FormatNames(path string) error {
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
		if header.IsReprocessed() {
			act = "R"
		}
		newFilename := fmt.Sprintf(nameFormat, header.GetHeadquarter(), header.GetStatementId(), header.GetPeriodInit().Format(printDateFormat),
			header.GetPeriodEnd().Format(printDateFormat), act, header.GetProcessingDate().Format(printDateFormat), header.GetLayoutVersion())
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

func (s Service) GetPeriodMap(path string) (map[time.Time]int, error) {
	dMap := make(map[time.Time]int)
	files, err := s.fileManager.GetFiles(path)
	if err != nil {
		return dMap, err
	}
	for _, f := range files {
		hData, err := s.GetHeaderData(path, f)
		if err != nil {
			continue
		}
		ds, err := hData.GetPeriodDates()
		if err != nil {
			continue
		}
		for _, d := range ds {
			if val, ok := dMap[d]; ok {
				dMap[d] = val + 1
			} else {
				dMap[d] = 1
			}
		}
	}
	return dMap, err
}

func (s Service) GetPeriod(path string) ([]time.Time, error) {
	dMap, err := s.GetPeriodMap(path)
	if err != nil {
		return make([]time.Time, 0), err
	}
	dates := make([]time.Time, 0)
	for d := range dMap {
		dates = append(dates, d)
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].After(dates[j]) })
	return dates, nil
}

func (s Service) GetPeriodGap(path string, initDate time.Time, endDate time.Time) ([]time.Time, error) {
	searchPeriod := make([]time.Time, 0)
	if initDate.Equal(time.Time{}) || endDate.Equal(time.Time{}) {
		return searchPeriod, fmt.Errorf("period is empty")
	}
	if initDate.After(endDate) {
		return searchPeriod, fmt.Errorf("initDate after endDate")
	}
	for t := initDate; !t.After(endDate); t = t.Add(24 * time.Hour) {
		searchPeriod = append(searchPeriod, t)
	}
	mdMap, err := s.GetPeriodMap(path)
	if err != nil {
		return make([]time.Time, 0), err
	}
	gaps := make([]time.Time, 0)
	for _, d := range searchPeriod {
		if _, ok := mdMap[d]; !ok {
			gaps = append(gaps, d)
		}
	}
	return gaps, nil
}
