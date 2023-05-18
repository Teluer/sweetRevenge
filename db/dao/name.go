package dao

import "sweetRevenge/db/dto"

func GetLeastUsedFirstName() string {
	var name dto.FirstName
	dao.db.Order("used_times asc, rand()").First(&name)
	name.UsedTimes++
	dao.db.Save(&name)
	return name.FirstName
}

func GetLeastUsedLastName() string {
	var name dto.LastName
	dao.db.Order("used_times asc, rand()").First(&name)
	name.UsedTimes++
	dao.db.Save(&name)
	return name.LastName
}
