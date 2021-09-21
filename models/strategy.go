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
	UsdtBalance float64 `json:"usdt_balance"` // usdt余额
	BtcBalance  float64 `json:"btc_balance"`  // btc余额
	EthBalance  float64 `json:"eth_balance"`  // eth余额
	RmbValue    float64 `json:"rmb_value"`    // 当前人民币估值
	PerMoney    float64 `json:"per_money"`    // 每次购买多少金额
}

func (s *Strategy) Save() error {
	return db.DB.Save(s).Error
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
