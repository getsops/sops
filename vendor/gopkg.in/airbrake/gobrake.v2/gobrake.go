package gobrake

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	SetLogger(log.New(os.Stderr, "gobrake: ", log.LstdFlags))
}

func SetLogger(l *log.Logger) {
	logger = l
}
