package main

import (
	"github.com/pkg/errors"
	"log"
	"os"

	"github.com/lavinas/cielo-edi/internal/handlers"
	"github.com/lavinas/cielo-edi/internal/utils/logger"
)

func main() {
	lg := logger.NewLogger()
	cm := handlers.NewCommandLine(lg)
	if err := cm.Run(os.Args); err != nil {
		err = errors.Wrap(err, "Error")
		log.Println(err)
	}
}
