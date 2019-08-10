package five_up_down

import (
	"encoding/json"
	"fmt"

	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/bitcoin/library/util"
	"gitlab.azbit.cn/web/bitcoin/models"
	"gitlab.azbit.cn/web/bitcoin/modules/huobi"
)

var logger = logging.MustGetLogger("modules/strategy/five_up_down")

type StrategyProcess struct {
	Strategy models.Strategy
	OrderMap map[string]*models.Order // 价格对应订单
}

var sp *StrategyProcess

const (
	STRATEGY_FIVE = 5
)

func Init(strategy models.Strategy) {
	// 查询历史订单
	ol, err := models.GetOrdersByStatus(strategy.ID, models.OrderStatusSuccess)
	if err != nil {
		logger.Error("five_up_down get order by status err:", err)
	}
	logger.Info("five_up_down history order num:", len(ol))
	om := make(map[string]*models.Order)
	for _, o := range ol {
		key := genOrderKey(strategy.ID, o.ExternalID)
		om[key] = o
	}
	sp = &StrategyProcess{
		Strategy: strategy,
		OrderMap: om,
	}
}

func StrategyDeal(kld *models.KLineData) {
	kld.Open = 11811.10
	logger.Info("strategy:", sp.Strategy.Name, "price:", kld.Open, " timestamp:", kld.Ts, " come in deal kline")
	if nil == sp {
		logger.Error("sp is nil...")
	}
	err := strategyProcessDeal(sp, kld)
	if err != nil {
		logger.Error("strategy process deal err:", err)
	}
}

// TODO：需要支持多线程
func strategyProcessDeal(sp *StrategyProcess, kld *models.KLineData) error {
	d, r := util.GetBackNum(int(kld.Open), sp.Strategy.Depth)
	logger.Info("d:", d, ", r:", r)
	rp := float64(d * sp.Strategy.Depth)
	var err error
	if r < STRATEGY_FIVE {
		// 如果没有买入则买入
		o := findBuyOrderByPrice(rp)
		if nil == o {
			o, err = order(sp.Strategy.FloatRate*rp, rp, models.OrderTypeBuy, kld.Ts)
			if err == nil {
				key := genOrderKey(sp.Strategy.ID, o.ExternalID)
				sp.OrderMap[key] = o
			}
		}
	} else {
		ol := findCanSaleOrders(rp)
		for _, bo := range ol {
			so, err := order(bo.Amount, bo.RefrencePrice+float64(STRATEGY_FIVE*sp.Strategy.Depth), models.OrderTypeSale, kld.Ts)
			if err == nil {
				err = settle(bo, so)
				if err == nil {
					key := genOrderKey(sp.Strategy.ID, bo.ExternalID)
					delete(sp.OrderMap, key)
				}
			}
		}
	}
	return err
}

func findBuyOrderByPrice(price float64) *models.Order {
	for _, o := range sp.OrderMap {
		if o.RefrencePrice == price && o.Type == models.OrderTypeBuy {
			return o
		}
	}
	return nil
}

func findCanSaleOrders(price float64) []*models.Order {
	ol := make([]*models.Order, 0)
	p := price - float64(5*sp.Strategy.Depth)
	for _, order := range sp.OrderMap {
		if order.RefrencePrice <= p && order.Type == models.OrderTypeBuy {
			ol = append(ol, order)
		}
		if order.Type == models.OrderTypeSale {
			logger.Error("order map should not have sale order")
		}
	}
	return ol
}

func genOrderKey(sid uint, externalID string) string {
	return fmt.Sprintf("%v_%v", sid, externalID)
}

func order(amount, price float64, orderType int, ts int64) (*models.Order, error) {
	logger.Info("order amount:", amount, "order type:", orderType, "order ts:", ts)
	o, err := huobi.HuobiPlaceOrder(sp.Strategy, "btcusdt", orderType, amount)
	if err != nil {
		return o, err
	}
	o.RefrencePrice = price
	o.Ts = ts
	err = o.Save()
	return o, err
}

// 结算买单和卖单
func settle(bo, so *models.Order) error {
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
		StrategyID:  sp.Strategy.ID,
		TotalAmount: sp.Strategy.TotalAmount,
		Depth:       sp.Strategy.Depth,
		FloatRate:   sp.Strategy.FloatRate,
		Capital:     bo.Money,
		InCome:      so.Money,
		Fee:         fee,
		Profit:      so.Money - bo.Money,
		Reason:      "five_up_down",
		Day:         util.GetTodayDay(),
		Orders:      string(idsByte),
	}
	err = p.Save()
	return err
}
