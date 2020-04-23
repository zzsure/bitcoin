package models

import "bitcoin/library/db"

const (
	OrderTypeUnkown = iota
	OrderTypeBuy
	OrderTypeSale
)

const (
	OrderStatusUnkown = iota
	OrderStatusSuccess
	OrderStatusBuy
	OrderStatusCancel
	OrderStatusSettle
)

type Order struct {
	Model
	StrategyID    uint    `json:"strategy_id"`    // 使用的策略id
	Money         float64 `json:"money"`          // 金额
	Price         float64 `json:"price"`          // 下单的时候价格
	Amount        float64 `json:"amount"`         // 成交量
	Fee           float64 `json:"fee"`            // 手续费，换算成USDT
	Type          int     `json:"type"`           // 下单类型，1为买入，2为卖出
	Status        int     `json:"status"`         // 状态，1为成单，2为下单，3为撤单，4为已结算
	Ts            int64   `json:"ts"`             // 下单时候K线时间戳，秒级别
	ExternalID    string  `json:"external_id"`    // 第三方下单的id
	RefrencePrice float64 `json:"refrence_price"` // k线参考价格
}

func (o *Order) Save() error {
	return db.DB.Save(o).Error
}

func GetOrdersByStatus(sid uint, status int) ([]*Order, error) {
	var os []*Order
	err := db.DB.Where("strategy_id = ? and status = ?", sid, status).Order("ts asc").Find(&os).Error
	return os, err
}
