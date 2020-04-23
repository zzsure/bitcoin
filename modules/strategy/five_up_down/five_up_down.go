package five_up_down

import (
	"bitcoin/library/util"
	"bitcoin/models"
	"bitcoin/modules/order"
	"github.com/op/go-logging"
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
		key := util.GetOrderKey(strategy.ID, o.ExternalID)
		om[key] = o
	}
	sp = &StrategyProcess{
		Strategy: strategy,
		OrderMap: om,
	}
}

func StrategyDeal(kld *models.KLineData) {
	//kld.Open = 11811.10
	if nil == sp {
		logger.Error("sp is nil...")
		return
	}
	logger.Info("strategy:", sp.Strategy.Name, "price:", kld.Open, " timestamp:", kld.Ts, " come in deal kline")
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
			o, err = order.Order(sp.Strategy, sp.Strategy.FloatRate*rp, rp, models.OrderTypeBuy, kld.Ts)
			if err == nil {
				key := util.GetOrderKey(sp.Strategy.ID, o.ExternalID)
				sp.OrderMap[key] = o
			}
		}
	} else {
		ol := findCanSaleOrders(rp)
		for _, bo := range ol {
			so, err := order.Order(sp.Strategy, bo.Amount, bo.RefrencePrice+float64(STRATEGY_FIVE*sp.Strategy.Depth), models.OrderTypeSale, kld.Ts)
			if err == nil {
				err = order.Settle(sp.Strategy, bo, so, "five_up_down")
				if err == nil {
					key := util.GetOrderKey(sp.Strategy.ID, bo.ExternalID)
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
