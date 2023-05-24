package dao

import (
	dto2 "sweetRevenge/src/db/dto"
)

func (d *gormDao) GetLeastUsedFirstName() string {
	var name dto2.FirstName
	d.db.Order("used_times asc, rand()").First(&name)
	name.UsedTimes++
	d.db.Save(&name)
	return name.FirstName
}

func (d *gormDao) GetLeastUsedLastName() string {
	var name dto2.LastName
	d.db.Order("used_times asc, rand()").First(&name)
	name.UsedTimes++
	d.db.Save(&name)
	return name.LastName
}
