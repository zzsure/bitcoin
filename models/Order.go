package models

import "gitlab.azbit.cn/web/bitcoin/library/db"

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
)

type Order struct {
	Model
	Strategy string  `json:"strategy"` // 使用的策略
	Money    float64 `json:"money"`    // 金额
	Price    float64 `json:"price"`    // 下单的时候价格
	Amount   float64 `json:"amount"`   // 成交量
	Fee      float64 `json:"fee"`      // 手续费
	Type     int     `json:"type"`     // 下单类型，1为买入，2为卖出
	Status   int     `json:"status"`   // 状态，1为成单，2为下单，3为撤单
	Ts       int64   `json:"ts"`       // 下单时候K线时间戳
}

func (o *Order) Save() error {
	return db.DB.Save(o).Error
}
