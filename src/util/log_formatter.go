package util

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type LogFormatter struct {
	log.Formatter
	Loc *time.Location
}

func (u LogFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.In(u.Loc)
	return u.Formatter.Format(e)
}
