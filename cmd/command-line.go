package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lavinas/cielo-edi/internal/core/domain"
	"github.com/lavinas/cielo-edi/internal/core/services"
	"github.com/lavinas/cielo-edi/internal/utils/file_manager"
	"github.com/lavinas/cielo-edi/internal/utils/string_parser"
)

var (
	funcMap = map[string]interface{}{
		"rename": rename,
		"gaps":   gaps,
	}
)

func rename(path string, args []string) error {
	parser := string_parser.NewStringParser()
	manager := file_manager.NewFileManager()
	header := domain.NewHeader(parser)
	service := services.NewService(manager, header)
	err := service.FormatNames(path)
	if err != nil {
		return err
	}
	return nil
}

func gapsExtraParam(args []string) (time.Time, time.Time, error) {
	zeroTime := time.Time{}
	if len(args) < 4 {
		err := fmt.Errorf("not enouth parameters (should by  ./command-line command path initialDate finalDate)")
		return zeroTime, zeroTime, err
	}
	initDate := args[3]
	finalDate := args[4]
	dInit, err := time.Parse("02/01/2006", initDate)
	if err != nil {
		err := fmt.Errorf("initial date error %v", err)
		return zeroTime, zeroTime, err
	}
	dFinal, err := time.Parse("02/01/2006", finalDate)
	if err != nil {
		err := fmt.Errorf("final date error %v", err)
		return zeroTime, zeroTime, err
	}
	return dInit, dFinal, nil
}

func gaps(path string, args []string) error {
	initDate, endDate, err := gapsExtraParam(args)
	if err != nil {
		return err
	}
	parser := string_parser.NewStringParser()
	manager := file_manager.NewFileManager()
	header := domain.NewHeader(parser)
	service := services.NewService(manager, header)
	dates, err := service.GetPeriodGap(path, initDate, endDate)
	if err != nil {
		return err
	}
	for _, date := range dates {
		log.Printf("%v\n", date)
	}

	return nil
}

func getArgs(args []string) (interface{}, string, error) {
	command := strings.ToLower(args[1])
	if command == "" {
		return nil, "", fmt.Errorf("command not found (should be ./command-line command path")
	}
	path := args[2]
	if path == "" {
		return nil, "", fmt.Errorf("command not found (should be ./command-line command path")
	}
	dir, err := os.Stat(path)
	if err != nil {
		return nil, "", err
	}
	if _, ok := funcMap[command]; !ok {
		return nil, "", fmt.Errorf("command %s not found (should be rename or gaps)", command)
	}
	if !dir.IsDir() {
		return nil, "", fmt.Errorf("dir %s do not exists", path)
	}
	return funcMap[command], path, nil
}

func exec(args []string) error {
	f, p, err := getArgs(args)
	if err != nil {
		return err
	}
	if err := f.(func(string, []string) error)(p, args); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := exec(os.Args); err != nil {
		log.Panic(err)
	}
}
