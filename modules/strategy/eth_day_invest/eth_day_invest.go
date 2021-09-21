package eth_day_invest

import (
	"bitcoin/library/util"
	"bitcoin/models"
	"bitcoin/modules/order"
	"github.com/op/go-logging"
	"time"
    "errors"
)

type StrategyProcess struct {
	Strategy models.Strategy
	DateMap  map[string]*models.Order // 日期订单
}

var logger = logging.MustGetLogger("modules/strategy/eth_day_invest")
var sp *StrategyProcess

func Init(strategy models.Strategy) {
	// 查询历史订单
	ol, err := models.GetOrdersByStatus(strategy.ID, models.OrderStatusSuccess)
	if err != nil {
		logger.Error("eth_day_invest get order by status err:", err)
	}
	logger.Info("eth_day_invest history order num:", len(ol))
	sp = &StrategyProcess{
		Strategy: strategy,
		DateMap:  make(map[string]*models.Order),
	}
	for _, o := range ol {
		addOrderToProcess(o)
	}
}

func StrategyDeal(kld *models.KLineData) {
	// kld.Open = 11811.10
	if nil == sp {
		//logger.Error("sp is nil...")
		return
	}
	logger.Info("eth strategy:", sp.Strategy.Name, "price:", kld.Open, " timestamp:", kld.Ts, " come in deal kline")
	err := strategyProcessDeal(kld)
	if err != nil {
		logger.Error("eth strategy process deal err:", err)
	}
}

func strategyProcessDeal(kld *models.KLineData) error {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(loc)

    // 最多定投3650天
    maxDays := 3650
    if len(sp.DateMap) > maxDays {
        return errors.New("max days: 1500")
    }

	// TODO：加锁，上一个不完成，不能下第二单
	date := util.GetDateByTime(t)
	// 8点开始投
	if 8 == util.GetCurlHour(t) {
		if _, ok := sp.DateMap[date]; !ok {
			logger.Info("order date: ", date)
			sp.DateMap[date] = new(models.Order)

			// 定投金额20$
			amount := 20.0

			o, err := order.EthOrder(sp.Strategy, amount, kld.Open, models.OrderTypeBuy, kld.Ts)
			if err != nil {
				return err
			}
			addOrderToProcess(o)
		} else {
			logger.Info("already buy")
			//return errors.New("")
		}
	}

	return nil
}

func addOrderToProcess(o *models.Order) {
	t := util.GetTimeByUnixTime(o.Ts)
	date := util.GetDateByTime(t)
	logger.Info("add order to process date: ", date)
	sp.DateMap[date] = o
}
