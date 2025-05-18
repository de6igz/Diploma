package config

import (
	"log"
	"time"
)

type AletheiaConfig struct {
	ProjectId        string
	CollectorAddress string
	ServiceName      string
	Environment      string
	Version          string
	AuthToken        string
	UseTLS           bool
	CaCertPath       string
	RPCTimeout       time.Duration
	Logger           Logger
}

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type defaultLogger struct{}

func (d *defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[SDK INFO] "+format, args...)
}

func (d *defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[SDK ERROR] "+format, args...)
}

func NewDefaultLogger() Logger {
	return &defaultLogger{}
}
