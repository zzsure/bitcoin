package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/bitcoin/conf"
	"gitlab.azbit.cn/web/bitcoin/library/util"
	"gitlab.azbit.cn/web/bitcoin/models"
)

var logger = logging.MustGetLogger("modules/socket")

var orderList []*models.Order
var totalAmount float64
var depth int
var lastSettle int64 // 上次结算时间

/*const (
    StrategyTypeNoChange = iota
    StrategyTypeDownBuy  // 下跌买入
    StrategyTypeSaleSettle   // 卖出盈利结算
    StrategyTypeSaleStopLoss // 卖出止损
)*/

func Init() {
	logger.Info("get all klds begin...")
	klds, err := models.GetKLineDataByTime(conf.Config.KLineData.Symbol, conf.Config.KLineData.From, conf.Config.KLineData.To)
	logger.Info("get all klds end...")
	if err != nil {
		logger.Error("kline data : ", err)
		return
	}
	for d := 1; d <= 10; d++ {
		conf.Config.Strategy.Floating.Depth = d
		for r := 0.01; r <= 0.30; r += 0.01 {
			conf.Config.Strategy.Floating.FloatRate = r
			depth = 0
			orderList = make([]*models.Order, 0)
			lastSettle = 0
			totalAmount = conf.Config.Strategy.Floating.TotalAmount
			start(klds, d, r)
		}
	}
}

func start(klds []*models.KLineData, d int, r float64) {
	logger.Info("d is : ", d, " r is: ", r)
	for idx, kld := range klds {
		if totalAmount <= 0.0 {
			logger.Error("blowing up...")
			return
		}
		if depth == 0 && len(orderList) == 0 {
			if kld.Ts-lastSettle < conf.Config.Strategy.Floating.Interval {
				//logger.Info("just wait interval k line buy...")
				continue
			}
			//orderList := make([]*models.Order, conf.Config.Strategy.Floating.Depth)
			orderList = make([]*models.Order, 0)
			err := order(kld.Open, models.OrderTypeBuy, kld.Ts)
			if err != nil {
				logger.Error("order fail: ", err)
			}
		}
		err := strategy(kld)
		if err != nil {
			logger.Error("strategy fail: ", err)
		}
		if idx == (len(klds)-1) && len(orderList) > 0 {
			i := len(orderList) - 1
			o := orderList[i]
			if models.OrderTypeBuy == o.Type {
				err := order(kld.High, models.OrderTypeSale, kld.Ts)
				if err != nil {
					logger.Error("order sale fail: ", err)
				}
				settle()
			}
		}
	}
}

func strategy(kld *models.KLineData) error {
	if len(orderList) == 0 {
		return errors.New("no order can use strategy")
	}
	//logger.Info("strategy k line...high price:", kld.High, " low price:", kld.Low)
	idx := len(orderList) - 1
	o := orderList[idx]
	lowPrice := (1.0 - conf.Config.Strategy.Floating.FloatRate) * o.Price
	if models.OrderTypeBuy == o.Type {
		expectIncome := o.Amount * kld.High
		for _, o := range orderList {
			if models.OrderTypeBuy == o.Type {
				expectIncome -= o.Money
			} else if models.OrderTypeSale == o.Type {
				expectIncome += o.Money
			}
		}
		if expectIncome > o.Money*conf.Config.Strategy.Floating.FloatRate {
			// 卖出盈利结算，复位
			err := order(kld.High, models.OrderTypeSale, kld.Ts)
			if err != nil {
				return err
			}
			return settle()
		}
		if depth < (conf.Config.Strategy.Floating.Depth-1) && kld.Low <= lowPrice {
			// 卖出止损，depth+1，缓存interval再买入
			err := order(lowPrice, models.OrderTypeSale, kld.Ts)
			if err != nil {
				return err
			}
			if depth >= conf.Config.Strategy.Floating.Depth-1 {
				settle()
				return errors.New("no enough fund")
			}
			depth += 1
		}
	} else if models.OrderTypeSale == o.Type {
		if kld.Ts-o.Ts >= conf.Config.Strategy.Floating.Interval {
			err := order(kld.Open, models.OrderTypeBuy, kld.Ts)
			if err != nil {
				return err
			}
		} else {
			//logger.Info("just wait strategy k line buy...")
		}
	}
	return nil
}

func settle() error {
	reason := fmt.Sprintf("depth%d", depth)
	if depth >= conf.Config.Strategy.Floating.Depth {
		reason = "nofund"
	}
	capital := 0.0
	income := 0.0
	fee := 0.0
	ids := make([]uint, len(orderList))
	for idx, o := range orderList {
		if models.OrderTypeBuy == o.Type {
			capital += o.Money
		} else if models.OrderTypeSale == o.Type {
			income += o.Money
		}
		fee += o.Fee
		lastSettle = o.Ts
		ids[idx] = o.ID
	}
	idsByte, _ := json.Marshal(ids)
	depth = 0
	//TODO:不太合理
	orderList = make([]*models.Order, 0)
	p := &models.Profit{
		StrategyID:  1,
		TotalAmount: conf.Config.Strategy.Floating.TotalAmount,
		Depth:       conf.Config.Strategy.Floating.Depth,
		FloatRate:   conf.Config.Strategy.Floating.FloatRate,
		Capital:     capital,
		InCome:      income,
		Fee:         fee,
		Profit:      income - capital,
		Reason:      reason,
		Day:         util.GetTodayDay(),
		Orders:      string(idsByte),
	}
	totalAmount = totalAmount + income - capital
	if totalAmount > conf.Config.Strategy.Floating.TotalAmount {
		totalAmount = conf.Config.Strategy.Floating.TotalAmount
	}
	err := p.Save()
	return err
}

func order(price float64, orderType int, ts int64) error {
	logger.Info("current depth: ", depth, "price: ", price, " ts: ", ts, "and type: ", orderType)
	money := 0.0
	fee := 0.0
	amount := 0.0
	if models.OrderTypeSale == orderType {
		if len(orderList) == 0 {
			return errors.New("no order can sale")
		}
		idx := len(orderList) - 1
		o := orderList[idx]
		if models.OrderTypeBuy != o.Type {
			return errors.New("last order is not buy")
		}
		fee = price * o.Amount * conf.Config.Huobi.SaleRates
		money = price*o.Amount - fee
		amount = o.Amount
	} else if models.OrderTypeBuy == orderType {
		per := totalAmount / (math.Pow(2, float64(conf.Config.Strategy.Floating.Depth)) - 1.0)
		money = math.Pow(2, float64(depth)) * per
		// 最后一次把余额都拿出来
		if depth == conf.Config.Strategy.Floating.Depth-1 {
			money = totalAmount
			for _, o := range orderList {
				if models.OrderTypeBuy == o.Type {
					money -= o.Money
				}
			}
		}
		amount = money / price
		fee = price * amount * conf.Config.Huobi.BuyRates
		amount -= amount * conf.Config.Huobi.BuyRates
	}
	logger.Info("current money: ", money, "amount: ", amount, " fee: ", fee)
	// 模拟下单即买入
	o := &models.Order{
		StrategyID: 1,
		Money:      money,
		Price:      price,
		Amount:     amount,
		Fee:        fee,
		Type:       orderType,
		Status:     models.OrderStatusSuccess, // 模拟下单即买入
		Ts:         ts,
	}
	err := o.Save()
	orderList = append(orderList, o)
	return err
}
