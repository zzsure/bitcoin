package day_float

import (
	"bitcoin/library/util"
	"bitcoin/models"
	"bitcoin/modules/order"
	"github.com/op/go-logging"
	"time"
)

type StrategyProcess struct {
	Strategy models.Strategy
	OrderMap map[string]*models.Order // 价格对应订单
	DateMap  map[string]*models.Order // 日期订单
}

var logger = logging.MustGetLogger("modules/strategy/day_float")
var sp *StrategyProcess

func Init(strategy models.Strategy) {
	// 查询历史订单
	ol, err := models.GetOrdersByStatus(strategy.ID, models.OrderStatusSuccess)
	if err != nil {
		logger.Error("day_float get order by status err:", err)
	}
	logger.Info("day_float history order num:", len(ol))
	sp = &StrategyProcess{
		Strategy: strategy,
		OrderMap: make(map[string]*models.Order),
		DateMap:  make(map[string]*models.Order),
	}
	for _, o := range ol {
		addOrderToProcess(sp.Strategy, o)
	}
}

func StrategyDeal(kld *models.KLineData) {
	// kld.Open = 11811.10
	logger.Info("strategy:", sp.Strategy.Name, "price:", kld.Open, " timestamp:", kld.Ts, " come in deal kline")
	if nil == sp {
		logger.Error("sp is nil...")
		return
	}
	err := strategyProcessDeal(kld)
	if err != nil {
		logger.Error("strategy process deal err:", err)
	}
}

func strategyProcessDeal(kld *models.KLineData) error {
	date := util.GetDateByTime(time.Now())
	if _, ok := sp.DateMap[date]; !ok {
		amount := sp.Strategy.PerMoney / kld.Open
		o, err := order.Order(sp.Strategy, amount, kld.Open, models.OrderTypeBuy, kld.Ts)
		if err == nil {
			addOrderToProcess(sp.Strategy, o)
		}
	}
	for _, bo := range sp.OrderMap {
		if (kld.Open-bo.Price)/bo.Price > sp.Strategy.FloatRate {
			so, err := order.Order(sp.Strategy, bo.Amount, kld.Open, models.OrderTypeSale, kld.Ts)
			if err == nil {
				err = order.Settle(sp.Strategy, bo, so, "day_float")
				if err == nil {
					removeOrderFromProcess(sp.Strategy, bo)
				}
			}
		}
	}
	return nil
}

func addOrderToProcess(s models.Strategy, o *models.Order) {
	key := util.GetOrderKey(s.ID, o.ExternalID)
	sp.OrderMap[key] = o
	t := util.GetTimeByUnixTime(o.Ts)
	date := util.GetDateByTime(t)
	sp.DateMap[date] = o
}

func removeOrderFromProcess(s models.Strategy, o *models.Order) {
	key := util.GetOrderKey(sp.Strategy.ID, o.ExternalID)
	delete(sp.OrderMap, key)
}
