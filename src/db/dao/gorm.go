package dao

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sweetRevenge/src/db/dto"
)

type GormDao struct {
	db *gorm.DB
}

var dao = open()

func open() *GormDao {
	log.Info("Opening connection to DB")
	dsn := "goblin:password1!@tcp(host.docker.internal:3306)/sweet?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to DB")
		panic(err)
	}
	return &GormDao{db}
}

func AutoMigrateAll() {
	dao.db.AutoMigrate(
		&dto.FirstName{},
		&dto.LastName{},
		&dto.Lady{},
		&dto.ManualOrder{},
		&dto.OrderHistory{})
}

func Insert(obj any) {
	log.WithField("obj", obj).Debug("Inserting data")
	dao.db.Create(obj)
}

func Delete(obj any) {
	log.WithField("obj", obj).Debug("Deleting data")
	dao.db.Where("1 = 1").Delete(obj)
}

func IsTableEmpty(obj any) bool {
	return dao.db.Limit(1).Find(obj).RowsAffected == 0
}

func FindFirst(obj any) {
	dao.db.Limit(1).Find(obj)
}
