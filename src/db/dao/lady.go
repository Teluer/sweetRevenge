package dao

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/src/db/dto"
)

// deprecated
func SaveNewLadies1(ladies []dto.Lady) {
	//remove all old records to avoid outdated phones
	if len(ladies) > 0 && !IsTableEmpty(&dto.Lady{}) {
		Delete(&dto.Lady{})
		Insert(&ladies)
		return
	}
}

func SaveNewLadies(ladies []dto.Lady) {
	log.Info("Saving new ladies")
	if IsTableEmpty(&dto.Lady{}) {
		log.Info("Ladies table is empty, populating")
		Insert(&ladies)
		return
	}
	log.Info("Ladies table has values, adding new, deleting outdated")

	phones := SelectPhones()
	var newLadies []dto.Lady
	var outdatedLadies []dto.Lady

NEW_LOOP:
	for _, lady := range ladies {
		for _, phone := range phones {
			if phone == lady.Phone {
				continue NEW_LOOP
			}
		}
		newLadies = append(newLadies, lady)
	}

OUTDATED_LOOP:
	for _, phone := range phones {
		for _, lady := range ladies {
			if phone == lady.Phone {
				continue OUTDATED_LOOP
			}
		}
		outdatedLadies = append(outdatedLadies, dto.Lady{Phone: phone})
	}

	log.Info(fmt.Sprintf("Inserting %d ladies", len(newLadies)))
	Insert(&newLadies)
	//log.Info(fmt.Sprintf("Deleting %d ladies", len(outdatedLadies)))
	//Delete(&outdatedLadies)
}

func SelectPhones() []string {
	var result []string
	dao.db.Model(&dto.Lady{}).Pluck("phone", &result)
	return result
}

func GetLeastUsedPhone() string {
	var lady dto.Lady
	dao.db.Order("used_times asc").First(&lady)
	lady.UsedTimes++
	dao.db.Save(&lady)
	return lady.Phone
}
