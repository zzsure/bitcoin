package strategy

import (
	"bitcoin/models"
	"bitcoin/modules/strategy/day_float"
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
		}
	}
}

func StrategyDeal(kld *models.KLineData) {
	floating.StrategyDeal(kld)
	five_up_down.StrategyDeal(kld)
	day_float.StrategyDeal(kld)
}
