package dao

import (
	dto2 "sweetRevenge/src/db/dto"
)

func GetLeastUsedFirstName() string {
	var name dto2.FirstName
	dao.db.Order("used_times asc, rand()").First(&name)
	name.UsedTimes++
	dao.db.Save(&name)
	return name.FirstName
}

func GetLeastUsedLastName() string {
	var name dto2.LastName
	dao.db.Order("used_times asc, rand()").First(&name)
	name.UsedTimes++
	dao.db.Save(&name)
	return name.LastName
}
