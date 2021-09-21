package strategy

import (
	"bitcoin/models"
	"bitcoin/modules/strategy/auto_invest"
	"bitcoin/modules/strategy/day_invest"
	"bitcoin/modules/strategy/day_float"
	"bitcoin/modules/strategy/eth_day_invest"
	"bitcoin/modules/strategy/five_up_down"
	"bitcoin/modules/strategy/floating"
	"github.com/op/go-logging"
	//"bitcoin/modules/strategy/history"
)

var logger = logging.MustGetLogger("modules/strategy")

func Init() {
	//history.Init()
	strategys, err := models.GetAllStrategys()
	if err != nil {
		logger.Error("get all strategy error...", err)
	}
	for _, s := range strategys {
		if s.Name == "floating" {
			floating.Init(s)
		} else if s.Name == "five_up_down" {
			five_up_down.Init(s)
		} else if s.Name == "day_float" {
			day_float.Init(s)
		} else if s.Name == "auto_invest" {
			auto_invest.Init(s)
		} else if s.Name == "day_invest" {
            day_invest.Init(s)
        } else if s.Name == "eth_day_invest" {
        	eth_day_invest.Init(s)
		}
	}
}

func StrategyDeal(kld *models.KLineData) {
	floating.StrategyDeal(kld)
	five_up_down.StrategyDeal(kld)
	day_float.StrategyDeal(kld)
	// 按金额定投，每周日7点定投，小于5000投50$，5000-10000：40$，10000-15000：30$，15000-20000：20$，大于20000：10$
	auto_invest.StrategyDeal(kld)
    // 按比特币数量定投，每天定投0.001
    day_invest.StrategyDeal(kld)
	// 按照20$每日定投ETH
	eth_day_invest.StrategyDeal(kld)
}
