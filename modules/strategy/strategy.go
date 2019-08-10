package strategy

import (
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/bitcoin/models"
	"gitlab.azbit.cn/web/bitcoin/modules/strategy/five_up_down"
	"gitlab.azbit.cn/web/bitcoin/modules/strategy/floating"
	//"gitlab.azbit.cn/web/bitcoin/modules/strategy/history"
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
		}
	}
}

func StrategyDeal(kld *models.KLineData) {
	floating.StrategyDeal(kld)
	five_up_down.StrategyDeal(kld)
}
