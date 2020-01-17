package models

import "bitcoin/library/db"

type Strategy struct {
	Model
	Name        string  `json:"name"`
	AccessKey   string  `json:"access_key"`
	SecretKey   string  `json:"secret_key"`
	AccountID   string  `json:"account_id"`
	TotalAmount float64 `json:"total_amount"` // 总金额仓位
	FloatRate   float64 `json:"float_rate"`   // 上下浮动的比例
	Depth       int     `json:"depth"`        // 最多下降和上升多少
	Interval    int64   `json:"interval"`     // 间隔多少s再次启用策略
	Status      int     `json:"status"`       // 0未启用，1启用
}

func GetAllStrategys() ([]Strategy, error) {
	var s []Strategy
	err := db.DB.Where("status = 1").Find(&s).Error
	return s, err
}

func GetStrategyByID(id int) (Strategy, error) {
	var s Strategy
	err := db.DB.Where("id = ? and status = 1", id).Find(&s).Error
	return s, err
}
