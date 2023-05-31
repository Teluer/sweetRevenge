package util

import (
	log "github.com/sirupsen/logrus"
)

func RecoverAndLog(flowName string) {
	err := recover()
	if err != nil {
		log.Error(flowName, " flow recovered from error: ", err)
	}
}

func RecoverAndLogAndDo(flowName string, action func()) {
	err := recover()
	if err != nil {
		log.Error(flowName, " flow recovered from error: ", err)
		action()
	}
}
