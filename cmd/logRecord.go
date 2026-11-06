package cmd

import (
	"time"
)

type LogRecord struct {
	Time time.Time
	Type string
	Text string
}
