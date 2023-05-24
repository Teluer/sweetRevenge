package dao

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sweetRevenge/src/db/dto"
)

type gormDao struct {
	db *gorm.DB
}

var Dao gormDao

func (d *gormDao) OpenDatabaseConnection() {
	log.Info("Opening connection to DB")
	dsn := "goblin:password1!@tcp(host.docker.internal:3306)/sweet?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to DB")
		panic(err)
	}
	d.db = db
}

func (d *gormDao) AutoMigrateAll() {
	d.db.AutoMigrate(
		&dto.FirstName{},
		&dto.LastName{},
		&dto.Lady{},
		&dto.OrderHistory{})
}

func (d *gormDao) Insert(obj any) {
	log.WithField("obj", obj).Debug("Inserting data")
	d.db.Create(obj)
}

func (d *gormDao) Delete(obj any) {
	log.WithField("obj", obj).Debug("Deleting data")
	d.db.Where("1 = 1").Delete(obj)
}

func (d *gormDao) IsTableEmpty(obj any) bool {
	return d.db.Limit(1).Find(obj).RowsAffected == 0
}
