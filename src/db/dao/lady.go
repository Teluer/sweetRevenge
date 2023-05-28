package dao

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/src/db/dto"
)

func (d *GormDao) SaveNewLadies(ladies []dto.Lady) {
	log.Info("Saving new ladies")
	if d.IsTableEmpty(&dto.Lady{}) {
		log.Info("Ladies table is empty, populating")
		d.Insert(&ladies)
		return
	}
	log.Info("Ladies table has values, adding new")

	phones := d.SelectPhones()
	var newLadies []dto.Lady

NEW_LOOP:
	for _, lady := range ladies {
		for _, phone := range phones {
			if phone == lady.Phone {
				continue NEW_LOOP
			}
		}
		newLadies = append(newLadies, lady)
	}

	if len(newLadies) > 0 {
		log.Info(fmt.Sprintf("Inserting %d ladies", len(newLadies)))
		d.Insert(&newLadies)
	} else {
		log.Info("No new ladies found")
	}
}

func (d *GormDao) SelectPhones() []string {
	var result []string
	d.db.Model(&dto.Lady{}).Pluck("phone", &result)
	return result
}

func (d *GormDao) GetLeastUsedPhone() string {
	var lady dto.Lady
	d.db.Order("used_times asc, rand()").First(&lady)
	lady.UsedTimes++
	d.db.Save(&lady)
	return lady.Phone
}
