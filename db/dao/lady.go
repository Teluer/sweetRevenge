package dao

import "sweetRevenge/db/dto"

func SaveNewLadies(ladies []dto.Lady) {
	if (IsTableEmpty(&dto.Lady{})) {
		Insert(&ladies)
		return
	}

	phones := SelectPhones()
	var newLadies []dto.Lady

MAIN_LOOP:
	for _, lady := range ladies {
		for _, phone := range phones {
			if phone == lady.Phone {
				continue MAIN_LOOP
			}
		}
		newLadies = append(newLadies, lady)
	}

	Insert(&newLadies)
}

func SelectPhones() []string {
	var result []string
	dao.db.Model(&dto.Lady{}).Pluck("phone", &result)
	return result
}

func GetLeastUsedPhone() string {
	var lady dto.Lady
	dao.db.Order("used_times asc, used_last").First(&lady)
	lady.UsedTimes++
	dao.db.Save(lady)
	return lady.Phone
}
