package models

import (
	//"github.com/op/go-logging"
	"time"

	"bitcoin/library/db"
)

type Model struct {
	ID        uint `gorm:"primary_key;auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//var logger = logging.MustGetLogger("model")

func CreateTable() {
	/*db.DB.DropTableIfExists(&KLineData{})
	db.DB.DropTableIfExists(&&Order{})
	db.DB.DropTableIfExists(&Profit{})*/
	_db := db.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	_db.LogMode(true)
	_db.CreateTable(&KLineData{})
	_db.CreateTable(&Order{})
	_db.CreateTable(&Profit{})
	_db.CreateTable(&Strategy{})
}

func MigrateTable() {
	_db := db.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	_db.LogMode(true)
	_db.AutoMigrate(&KLineData{})
	_db.AutoMigrate(&Order{})
	_db.AutoMigrate(&Profit{})
	_db.AutoMigrate(&Strategy{})
}
