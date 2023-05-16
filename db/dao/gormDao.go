package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormDao struct {
	db *gorm.DB
}

var dao = open()

func open() *GormDao {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "goblin:password1!@tcp(localhost:3306)/sweet?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &GormDao{db}
}

func Insert(obj any) {
	dao.db.Create(obj)
}

func IsTableEmpty(obj any) bool {
	return dao.db.Limit(1).Find(obj).RowsAffected == 0
}
