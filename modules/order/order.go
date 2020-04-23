package order

import (
	"bitcoin/library/util"
	"bitcoin/models"
	"bitcoin/modules/huobi"
	"encoding/json"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("modules/order")

func Order(s models.Strategy, amount, price float64, orderType int, ts int64) (*models.Order, error) {
	logger.Info("order amount:", amount, "order type:", orderType, "order ts:", ts)
	o, err := huobi.HuobiPlaceOrder(s, "btcusdt", orderType, amount)
	if err != nil {
		return o, err
	}
	o.RefrencePrice = price
	o.Ts = ts
	err = o.Save()
	return o, err
}

// 结算买单和卖单
func Settle(s models.Strategy, bo, so *models.Order, reason string) error {
	// TODO:改成事务
	bo.Status = models.OrderStatusSettle
	err := bo.Save()
	if err != nil {
		return err
	}
	so.Status = models.OrderStatusSettle
	err = so.Save()
	if err != nil {
		return err
	}
	fee := so.Fee + bo.Fee
	ids := []uint{bo.ID, so.ID}
	idsByte, _ := json.Marshal(ids)
	p := &models.Profit{
		StrategyID:  s.ID,
		TotalAmount: s.TotalAmount,
		Depth:       s.Depth,
		FloatRate:   s.FloatRate,
		Capital:     bo.Money,
		InCome:      so.Money,
		Fee:         fee,
		Profit:      so.Money - bo.Money,
		Reason:      reason,
		Day:         util.GetTodayDay(),
		Orders:      string(idsByte),
	}
	err = p.Save()
	return err
}
