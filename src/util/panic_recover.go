package util

import (
	log "github.com/sirupsen/logrus"
	"sweetRevenge/src/db/dao"
)

func RecoverAndLog(flowName string) {
	err := recover()
	if err != nil {
		log.Error(flowName, " flow recovered from error: ", err)
	}
}

func RecoverAndRollbackAndLog(flowName string, tx dao.Database) {
	err := recover()
	if err != nil {
		tx.RollbackTransaction()
		log.Error(flowName, " flow recovered from error: ", err)
	}
}
