package infra

import (
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "[search-bff] ", log.LstdFlags|log.Lshortfile)
}
