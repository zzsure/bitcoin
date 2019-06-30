package cron

import (
	"encoding/json"

	"github.com/op/go-logging"
	"github.com/robfig/cron"
	"gitlab.azbit.cn/web/bitcoin/models"
	"gitlab.azbit.cn/web/bitcoin/modules/huobi"
	"gitlab.azbit.cn/web/bitcoin/modules/strategy"
)

var c *cron.Cron
var lastKld *models.KLineData
var logger = logging.MustGetLogger("modules/cron")

func Init() {
	c = cron.New()
	getHuobiKLineCron()
	c.Start()
}

func getHuobiKLineCron() {
	c.AddFunc("*/10 * * * * *", func() {
		getHuobiKLine()
	})
}

func getHuobiKLine() {
	r := huobi.GetKLine("btcusdt", "1min", 1)
	if len(r.Data) > 0 {
		data := r.Data[0]
		kld := &models.KLineData{
			Kid:    data.ID,
			Amount: data.Amount,
			Count:  data.Count,
			Open:   data.Open,
			Close:  data.Close,
			Low:    data.Low,
			High:   data.High,
			Vol:    data.Vol,
			Ch:     "market.btcusdt.kline.1min",
			Ts:     data.ID,
		}
		info, _ := json.Marshal(kld)
		logger.Info("recv line: %s", string(info))
		strategy.StrategyDeal(kld)
		if lastKld != nil && lastKld.Ts != kld.Ts {
			err := kld.Save()
			if err != nil {
				logger.Error("save kline data to db: ", err)
			}
		}
		lastKld = kld
	}
}
