package floating

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/op/go-logging"
	"bitcoin/conf"
	"bitcoin/library/util"
	"bitcoin/models"
	"bitcoin/modules/huobi"
)

var logger = logging.MustGetLogger("modules/strategy/floating")

type StrategyProcess struct {
	Strategy    models.Strategy
	OrderList   []*models.Order
	TotalAmount float64
	Depth       int
	LastSettle  int64 // 上次结算时间
}

var spMap map[uint]*StrategyProcess

// TODO: 下单没成单的处理
func Init(strategy models.Strategy) {
	if nil == spMap {
		spMap = make(map[uint]*StrategyProcess)
	}
	// 查询历史订单
	ol, err := models.GetOrdersByStatus(strategy.ID, models.OrderStatusSuccess)
	if err != nil {
		logger.Error("get order by status err:", err)
	}
	logger.Info("history order num:", len(ol))
	ta := huobi.GetCurrencyBalance(strategy, "usdt")
	//ta = 30.0
	logger.Info("account:", strategy.AccountID, ", balance:", ta)
	d := 0
	for _, o := range ol {
		if o.Type == models.OrderTypeBuy {
			ta += o.Money
		} else if o.Type == models.OrderTypeSale {
			d++
		}
	}
	if ta > strategy.TotalAmount {
		ta = strategy.TotalAmount
	}
	// TODO:未成单的处理
	spMap[strategy.ID] = &StrategyProcess{
		Strategy:    strategy,
		OrderList:   ol,
		TotalAmount: ta,
		Depth:       d,
		LastSettle:  0,
	}
}

func StrategyDeal(kld *models.KLineData) {
	logger.Info("come in deal kline")
	for _, sp := range spMap {
		/*orderDetail := huobi.PlaceDetail(sp.Strategy, "38795364910")
		  info, _ := json.Marshal(orderDetail)
		  logger.Info("order detail:", string(info))
		  ta := huobi.GetCurrencyBalance(sp.Strategy, "btc")
		  logger.Info("account:", sp.Strategy.AccountID, ", balance:", ta)*/
		logger.Info("strategy:", sp.Strategy.Name, " timestamp:", kld.Ts)
		err := strategyProcessDeal(sp, kld)
		if err != nil {
			logger.Error("strategy process deal err:", err)
		}
	}
}

func strategyProcessDeal(sp *StrategyProcess, kld *models.KLineData) error {
	if sp.Depth == 0 && len(sp.OrderList) == 0 {
		if kld.Ts-sp.LastSettle < sp.Strategy.Interval {
			return errors.New("just wait interval k line buy...")
		}
		sp.OrderList = make([]*models.Order, 0)
		err := order(sp, kld.Open, models.OrderTypeBuy, kld.Ts)
		if err != nil {
			return err
		}
	}
	idx := len(sp.OrderList) - 1
	o := sp.OrderList[idx]
	lowPrice := (1.0 - sp.Strategy.FloatRate) * o.Price
	if models.OrderTypeBuy == o.Type {
		expectIncome := util.Float64Precision(o.Amount, 4, false) * kld.High
		logger.Info("expect income money is:", expectIncome)
		for _, o := range sp.OrderList {
			if models.OrderTypeBuy == o.Type {
				expectIncome -= o.Money
			} else if models.OrderTypeSale == o.Type {
				expectIncome += o.Money
			}
		}
		logger.Info("expect income money cal is:", expectIncome)
		logger.Info("expect o.money is:", o.Money*sp.Strategy.FloatRate)
		if expectIncome > o.Money*sp.Strategy.FloatRate {
			// 卖出盈利结算，复位
			err := order(sp, kld.High, models.OrderTypeSale, kld.Ts)
			if err != nil {
				return err
			}
			return settle(sp)
		}
		logger.Info("kld low price:", kld.Low, " unexpect low price:", lowPrice)
		if sp.Depth < (sp.Strategy.Depth-1) && kld.Low <= lowPrice {
			// 卖出止损，depth+1，缓存interval再买入
			err := order(sp, lowPrice, models.OrderTypeSale, kld.Ts)
			if err != nil {
				return err
			}
			if sp.Depth >= sp.Strategy.Depth-1 {
				settle(sp)
				return errors.New("no enough fund")
			}
			sp.Depth += 1
		}
	} else if models.OrderTypeSale == o.Type {
		if kld.Ts-o.Ts >= sp.Strategy.Interval {
			err := order(sp, kld.Open, models.OrderTypeBuy, kld.Ts)
			if err != nil {
				return err
			}
		} else {
			logger.Info("just wait strategy k line buy...")
		}
	}
	return nil
}

func order(sp *StrategyProcess, price float64, orderType int, ts int64) error {
	logger.Info("current depth: ", sp.Depth, "price: ", price, " ts: ", ts, "and type: ", orderType)
	money := 0.0
	fee := 0.0
	amount := 0.0
	if models.OrderTypeSale == orderType {
		if len(sp.OrderList) == 0 {
			return errors.New("no order can sale")
		}
		idx := len(sp.OrderList) - 1
		o := sp.OrderList[idx]
		if models.OrderTypeBuy != o.Type {
			return errors.New("last order is not buy")
		}
		fee = price * util.Float64Precision(o.Amount, 4, false) * conf.Config.Huobi.SaleRates
		money = price*util.Float64Precision(o.Amount, 4, false) - fee
		amount = util.Float64Precision(o.Amount, 4, false)
	} else if models.OrderTypeBuy == orderType {
		per := sp.TotalAmount / (math.Pow(2, float64(sp.Strategy.Depth)) - 1.0)
		money = math.Pow(2, float64(sp.Depth)) * per
		// 最后一次把余额都拿出来
		if sp.Depth == sp.Strategy.Depth-1 {
			money = huobi.GetCurrencyBalance(sp.Strategy, "usdt")
			if money > sp.Strategy.TotalAmount {
				money = sp.Strategy.TotalAmount
			}
		}
		amount = money / price
		fee = price * amount * conf.Config.Huobi.BuyRates
		amount -= amount * conf.Config.Huobi.BuyRates
	}
	logger.Info("current money: ", money, "amount: ", amount, " fee: ", fee)

	// 火币下单
	if models.OrderTypeBuy == orderType {
		amount = money
	}
	o, err := huobi.HuobiPlaceOrder(sp.Strategy, "btcusdt", orderType, amount)
	if err != nil {
		return err
	}
	o.RefrencePrice = price
	o.Ts = ts
	err = o.Save()

	sp.OrderList = append(sp.OrderList, o)
	return err
}

func settle(sp *StrategyProcess) error {
	reason := fmt.Sprintf("depth%d", sp.Depth)
	if sp.Depth >= sp.Strategy.Depth {
		reason = "nofund"
	}
	capital := 0.0
	income := 0.0
	fee := 0.0
	ids := make([]uint, len(sp.OrderList))
	// TODO:改成事务
	for idx, o := range sp.OrderList {
		if models.OrderTypeBuy == o.Type {
			capital += o.Money
		} else if models.OrderTypeSale == o.Type {
			income += o.Money
		}
		fee += o.Fee
		sp.LastSettle = o.Ts
		ids[idx] = o.ID
		o.Status = models.OrderStatusSettle
		err := o.Save()
		if err != nil {
			//logger.Error("settle order status err:", err)
			return err
		}
	}
	idsByte, _ := json.Marshal(ids)
	sp.Depth = 0
	//TODO:不太合理
	sp.OrderList = make([]*models.Order, 0)
	p := &models.Profit{
		StrategyID:  sp.Strategy.ID,
		TotalAmount: sp.Strategy.TotalAmount,
		Depth:       sp.Strategy.Depth,
		FloatRate:   sp.Strategy.FloatRate,
		Capital:     capital,
		InCome:      income,
		Fee:         fee,
		Profit:      income - capital,
		Reason:      reason,
		Day:         util.GetTodayDay(),
		Orders:      string(idsByte),
	}
	sp.TotalAmount = huobi.GetCurrencyBalance(sp.Strategy, "usdt")
	if sp.TotalAmount > sp.Strategy.TotalAmount {
		sp.TotalAmount = sp.Strategy.TotalAmount
	}
	err := p.Save()
	return err
}

/*func start(sp *StrategyProcess, klds []*models.KLineData, d int, r float64) {
	logger.Info("d is : ", d, " r is: ", r)
	for idx, kld := range klds {
		if sp.TotalAmount <= 0.0 {
			logger.Error("blowing up...")
			return
		}
		if sp.Depth == 0 && len(sp.OrderList) == 0 {
			if kld.Ts - sp.LastSettle < sp.Strategy.Interval {
				logger.Info("just wait interval k line buy...")
				continue
			}
			//orderList := make([]*models.Order, conf.Config.Strategy.Floating.Depth)
			sp.OrderList = make([]*models.Order, 0)
			err := order(sp, kld.Open, models.OrderTypeBuy, kld.Ts)
			if err != nil {
				logger.Error("order fail: ", err)
			}
		}
		err := strategyProcessDeal(sp, kld)
		if err != nil {
			logger.Error("strategy fail: ", err)
		}
		if idx == (len(klds)-1) && len(sp.OrderList) > 0 {
			i := len(sp.OrderList) - 1
			o := sp.OrderList[i]
			if models.OrderTypeBuy == o.Type {
				err := order(sp, kld.High, models.OrderTypeSale, kld.Ts)
				if err != nil {
					logger.Error("order sale fail: ", err)
				}
				settle(sp)
			}
		}
	}
}*/
