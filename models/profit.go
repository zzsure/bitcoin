package models

import "gitlab.azbit.cn/web/bitcoin/library/db"

type Profit struct {
	Model
	StrategyID  uint    `json:"stragety_id"`  // 策略
	TotalAmount float64 `json:"total_amount"` // 总资产
	Depth       int     `json:"depth"`        // 最多下跌几次
	FloatRate   float64 `json:"float_rate"`   // 上下浮动的比例
	Capital     float64 `json:"capital"`      // 投入
	InCome      float64 `json:"in_come"`      // 收入
	Fee         float64 `json:"fee"`          // 手续费
	Profit      float64 `json:"profit"`       // 利润
	Reason      string  `json:"reason"`       // 为什么结算，depth1...depth8, nofund
	Day         string  `json:"day"`          // 日期
	Orders      string  `json:"orders"`       // 订单id列表
}

func (p *Profit) Save() error {
	return db.DB.Save(p).Error
}
