package models

import (
	//"github.com/op/go-logging"
	"time"

	"gitlab.azbit.cn/web/bitcoin/library/db"
)

type Model struct {
	ID        uint `gorm:"primary_key;auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

//var logger = logging.MustGetLogger("model")

func CreateTable() {
	db.DB.DropTableIfExists(&KLineData{})
	_db := db.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	_db.LogMode(true)
	_db.CreateTable(&KLineData{})
}

func MigrateTable() {
	_db := db.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	_db.LogMode(true)
	_db.AutoMigrate(&KLineData{})
}
