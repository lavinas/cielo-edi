package main

import (
	"github.com/lavinas/cielo-edi/internal/core/domain"
	"github.com/lavinas/cielo-edi/internal/core/services/rename"
	"github.com/lavinas/cielo-edi/internal/utils/string_parser"
)

const (
	path string = "/home/paulo/Desktop/nespresso/arquivos-x"
)

func main() {
	string_parser := string_parser.NewStringParser()
	header := domain.NewHeader(string_parser)
	rename := rename.NewRenameService(*header)
	err := rename.FormatFilesName(path)
	if err != nil {
		panic(err)
	}
}
