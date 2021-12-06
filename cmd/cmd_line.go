package main

import (
	"github.com/pkg/errors"
	"log"
	"os"

	"github.com/lavinas/cielo-edi/internal/core/ports"
	"github.com/lavinas/cielo-edi/internal/handlers"
)

func main() {
	var cm ports.CommandLineInterface = handlers.NewCommandLine()
	if err := cm.Run(os.Args); err != nil {
		err = errors.Wrap(err, "Error")
		log.Println(err)
	}
}
