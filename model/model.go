package model

import (
	"github.com/op/go-logging"
	"time"
)

type Model struct {
	ID        uint `gorm:"primary_key;auto_increment"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var logger = logging.MustGetLogger("model")

