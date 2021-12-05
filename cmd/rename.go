package main

import (
	"github.com/lavinas/cielo-edi/internal/core/domain"
	"github.com/lavinas/cielo-edi/internal/core/services"
	"github.com/lavinas/cielo-edi/internal/utils/string_parser"
	"github.com/lavinas/cielo-edi/internal/utils/file_manager"
)

const (
	path string = "/home/paulo/Desktop/nespresso/arquivos1"
)

func rename() error {
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

func main() {
	err := rename()
	if err != nil {
		panic(err)
	}
	
}
