package dao

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sweetRevenge/src/db/dto"
)

type Database interface {
	OpenDatabaseConnection(string)
	AutoMigrateAll()
	Insert(obj any)
	Delete(obj any)
	IsTableEmpty(obj any) bool
	ValidateDataIntegrity() bool

	SelectPhones() []string
	GetLeastUsedPhone() string
	SaveNewLadies(ladies []dto.Lady)
	GetLeastUsedFirstName() string
	GetLeastUsedLastName() string
}

type GormDao struct {
	db *gorm.DB
}

var Dao Database = &GormDao{}

func (d *GormDao) OpenDatabaseConnection(dsn string) {
	log.Info("Opening connection to DB")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to DB")
		panic(err)
	}
	d.db = db
}

func (d *GormDao) AutoMigrateAll() {
	d.db.AutoMigrate(
		&dto.FirstName{},
		&dto.LastName{},
		&dto.Lady{},
		&dto.OrderHistory{})
}

func (d *GormDao) Insert(obj any) {
	log.WithField("obj", obj).Debug("Inserting data")
	d.db.Create(obj)
}

func (d *GormDao) Delete(obj any) {
	log.WithField("obj", obj).Debug("Deleting data")
	d.db.Where("1 = 1").Delete(obj)
}

func (d *GormDao) IsTableEmpty(obj any) bool {
	return d.db.Limit(1).Find(obj).RowsAffected == 0
}

func (d *GormDao) ValidateDataIntegrity() bool {
	return !(Dao.IsTableEmpty(&dto.FirstName{}) ||
		Dao.IsTableEmpty(&dto.LastName{}) ||
		Dao.IsTableEmpty(&dto.Lady{}))
}
