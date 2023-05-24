package dao

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/src/db/dto"
)

// deprecated
func (d *gormDao) SaveNewLadies1(ladies []dto.Lady) {
	//remove all old records to avoid outdated phones
	if len(ladies) > 0 && !d.IsTableEmpty(&dto.Lady{}) {
		d.Delete(&dto.Lady{})
		d.Insert(&ladies)
		return
	}
}

func (d *gormDao) SaveNewLadies(ladies []dto.Lady) {
	log.Info("Saving new ladies")
	if d.IsTableEmpty(&dto.Lady{}) {
		log.Info("Ladies table is empty, populating")
		d.Insert(&ladies)
		return
	}
	log.Info("Ladies table has values, adding new, deleting outdated")

	phones := d.SelectPhones()
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
	d.Insert(&newLadies)
	//log.Info(fmt.Sprintf("Deleting %d ladies", len(outdatedLadies)))
	//Delete(&outdatedLadies)
}

func (d *gormDao) SelectPhones() []string {
	var result []string
	d.db.Model(&dto.Lady{}).Pluck("phone", &result)
	return result
}

func (d *gormDao) GetLeastUsedPhone() string {
	var lady dto.Lady
	d.db.Order("used_times asc").First(&lady)
	lady.UsedTimes++
	d.db.Save(&lady)
	return lady.Phone
}
