package handlers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lavinas/cielo-edi/internal/core/domain"
	"github.com/lavinas/cielo-edi/internal/core/ports"
	"github.com/lavinas/cielo-edi/internal/core/services"
	"github.com/lavinas/cielo-edi/internal/utils/file_manager"
	"github.com/lavinas/cielo-edi/internal/utils/string_parser"
)

var (
	funcMap = map[string]interface{}{
		"rename":  rename,
		"gaps":    gaps,
		"periods": periods,
	}
	acquirerMap = map[string]ports.HeaderDataInterface{
		"cielovendas":       &domain.HeaderCielo{Statement: "vendas"},
		"cielofinanceiro":   &domain.HeaderCielo{Statement: "financeiro"},
		"cieloantecipacoes": &domain.HeaderCielo{Statement: "antecipacoes"},
		"cieloalelo":        &domain.HeaderCielo{Statement: "alelo"}, 
		"redecredito":       &domain.HeaderRedeCredit{Statement: "credito"},
		"rededebito":        &domain.HeaderRedeDebt{Statement: "debito"},
		"redefinanceiro":    &domain.HeaderRedeFin{Statement: "financeiro"},
		"getnet":            &domain.HeaderGetnet{},
	}
	parserTypeMap = map[string]string{
		"cielovendas":       "position",
		"cielofinanceiro":   "position",
		"cieloantecipacoes": "position",
		"cieloalelo":        "position",
		"redecredito":       "position",
		"rededebito":        "csv",
		"redefinanceiro":    "position",
		"getnet":            "position",
	}
)

type CommandLine struct {
	logger ports.LoggerInterface
}

func NewCommandLine(logger ports.LoggerInterface) *CommandLine {
	return &CommandLine{logger: logger}
}

func (cm CommandLine) Run(args []string) error {
	function, headerData, parserType, path, err := getArgs(args)
	if err != nil {
		return err
	}
	parser := string_parser.NewStringParser(parserType)
	manager := file_manager.NewFileManager()
	header := domain.NewHeader(headerData, parser)
	service := services.NewService(manager, header)
	if err := function.(func(ports.LoggerInterface, ports.ServiceInterface, string, []string) error)(cm.logger, service, path, args); err != nil {
		return err
	}
	return nil
}

func rename(log ports.LoggerInterface, service ports.ServiceInterface, path string, args []string) error {
	logger, err := service.FormatNames(path)
	if err != nil {
		return err
	}
	for _, logLine := range logger {
		log.Println(logLine)
	}
	return nil
}

func gaps(log ports.LoggerInterface, service ports.ServiceInterface, path string, args []string) error {
	initDate, endDate, err := gapsExtraParam(args)
	if err != nil {
		return err
	}
	dates, err := service.GetGapGrouped(path, initDate, endDate)
	if err != nil {
		return err
	}
	for _, date := range dates {
		log.Printf("%s", date)
	}
	return nil
}

func periods(log ports.LoggerInterface, service ports.ServiceInterface, path string, args []string) error {
	dates, err := service.GetPeriodGrouped(path)
	if err != nil {
		return err
	}
	for _, date := range dates {
		log.Printf("%s", date)
	}
	return nil
}

func gapsExtraParam(args []string) (time.Time, time.Time, error) {
	zeroTime := time.Time{}
	if len(args) < 5 {
		err := fmt.Errorf("not enouth parameters (should by ./command-line command path initialDate finalDate)")
		return zeroTime, zeroTime, err
	}
	initDate := args[4]
	finalDate := args[5]
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

func getCommand(args []string) (interface{}, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("command not found (should be ./command-line command acquirer path")
	}
	command := strings.ToLower(args[1])
	if command == "" {
		return nil, fmt.Errorf("command not found (should be ./command-line command acquirer path")
	}
	if _, ok := funcMap[command]; !ok {
		return nil, fmt.Errorf("command %s not found (should be rename, gaps or periods)", command)
	}
	return funcMap[command], nil
}

func getAcquirer(args []string) (ports.HeaderDataInterface, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("command not found (should be ./command-line command acquirer path")
	}
	acquirer := args[2]
	if acquirer == "" {
		return nil, fmt.Errorf("command not found (should be ./command-line command acquirer path")
	}
	if _, ok := acquirerMap[acquirer]; !ok {
		return nil, fmt.Errorf("acquirer name %s not found (should be cielovendas, cielofinanceiro, cieloantecipacoes, redecredito, redebito, redefinanceiro, getnet)", acquirer)
	}
	return acquirerMap[acquirer], nil
}

func getParserType(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("command not found (should be ./command-line command acquirer path")
	}
	acquirer := args[2]
	if acquirer == "" {
		return "", fmt.Errorf("command not found (should be ./command-line command acquirer path")
	}
	if _, ok := acquirerMap[acquirer]; !ok {
		return "", fmt.Errorf("acquirer name %s not found (should be cielovendas, cielofinanceiro, cieloantecipacoes, redecredito, redebito, redefinanceiro, getnet)", acquirer)
	}
	return parserTypeMap[acquirer], nil
}

func getPath(args []string) (string, error) {
	if len(args) < 4 {
		return "", fmt.Errorf("wrong number of parameters (should be ./command-line command acquirer path")
	}
	path := args[3]
	if path == "" {
		return "", fmt.Errorf("command not found (should be ./command-line command acquirer path")
	}
	dir, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !dir.IsDir() {
		return "", fmt.Errorf("dir %s do not exists", path)
	}
	return path, nil
}

func getArgs(args []string) (interface{}, ports.HeaderDataInterface, string, string, error) {
	if len(args) < 4 {
		return nil, nil, "", "", fmt.Errorf("wrong number of parameters (should be ./command-line command acquirer path")
	}
	command, err := getCommand(args)
	if err != nil {
		return nil, nil, "", "", err
	}
	acquirer, err := getAcquirer(args)
	if err != nil {
		return nil, nil, "", "", err
	}
	path, err := getPath(args)
	if err != nil {
		return nil, nil, "", "", err
	}
	parseType, err := getParserType(args)
	if err != nil {
		return nil, nil, "", "", err
	}
	return command, acquirer, parseType, path, nil
}
