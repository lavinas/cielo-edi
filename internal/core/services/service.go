package services

import (
	"fmt"
	"io/fs"
	"sort"
	"time"

	"github.com/lavinas/cielo-edi/internal/core/ports"
)

const (
	nameFormat      string = "%s-%010d-%s-%s-%s-%s-%s-L%03d.txt"
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
		return nil, fmt.Errorf("error parsing")
	}
	if !s.header.IsValid() {
		return nil, fmt.Errorf("invalid file")
	}
	d := s.header.GetData()
	return d, nil
}

func (s Service) FormatNames(path string) ([]string, error) {
	logger := make([]string, 0)
	files, err := s.fileManager.GetFiles(path)
	if err != nil {
		return []string{}, err
	}
	for _, file := range files {
		h, err := s.GetHeaderData(path, file)
		if err != nil {
			logger = append(logger, fmt.Sprintf("No: %s - %v", file.Name(), err))
			continue
		}
		act := "N"
		if h.IsReprocessed() {
			act = "R"
		}
		newName := fmt.Sprintf(nameFormat, h.GetAcquirer(), h.GetHeadquarter(), h.GetStatementId(),
			h.GetPeriodInit().Format(printDateFormat), h.GetPeriodEnd().Format(printDateFormat), act,
			h.GetProcessingDate().Format(printDateFormat), h.GetLayoutVersion())
		err = s.fileManager.RenameFile(path, file.Name(), newName)
		if err != nil {
			logger = append(logger, fmt.Sprintf("No: %s - %v", file.Name(), err))
			continue
		}
		logger = append(logger, fmt.Sprintf("Yes: %s - %s", file.Name(), newName))
	}
	return logger, nil
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
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })
	return dates, nil
}

func (s Service) GetGap(path string, initDate time.Time, endDate time.Time) ([]time.Time, error) {
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

func (s Service) GetGrouped(dates []time.Time) []string {
	pInit := time.Time{}
	pEnd := time.Time{}
	ret := make([]string, 0)
	for _, date := range dates {
		if pInit.Equal(time.Time{}) {
			pInit = date
			pEnd = date
		}
		if date.After(pEnd.Add(24 * time.Hour)) {
			strRet := fmt.Sprintf("%s - %s", pInit.Format("02/01/2006"), pEnd.Format("02/01/2006"))
			ret = append(ret, strRet)
			pInit = date
			pEnd = date
		} else {
			pEnd = date
		}
	}
	if !pInit.Equal(time.Time{}) {
		strRet := fmt.Sprintf("%s - %s", pInit.Format("02/01/2006"), pEnd.Format("02/01/2006"))
		ret = append(ret, strRet)
	}
	return ret
}

func (s Service) GetGapGrouped(path string, initDate time.Time, endDate time.Time) ([]string, error) {
	dates, err := s.GetGap(path, initDate, endDate)
	if err != nil {
		return []string{}, err
	}
	return s.GetGrouped(dates), nil
}

func (s Service) GetPeriodGrouped(path string) ([]string, error) {
	dates, err := s.GetPeriod(path)
	if err != nil {
		return []string{}, err
	}
	return s.GetGrouped(dates), nil
}
