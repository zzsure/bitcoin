package cron

import (
	"bitcoin/conf"
	"bitcoin/consts"
	"bitcoin/library/redis"
	"encoding/json"
	"time"

	"github.com/op/go-logging"
	"github.com/robfig/cron"
	"bitcoin/models"
	"bitcoin/modules/huobi"
	"bitcoin/modules/strategy"
)

var c *cron.Cron
var lastKld *models.KLineData
var logger = logging.MustGetLogger("modules/cron")

func Init() {
	c = cron.New()
	getHuobiKLineCron()
	getBalanceCron()
	c.Start()
}

func getHuobiKLineCron() {
	c.AddFunc("@every 10s", func() {
		logger.Info("get huobi kline cron begin")
		getHuobiKLine()
		logger.Info("get huobi kline cron end")
	})
}

func getBalanceCron() {
	c.AddFunc("@every 1h", func() {
		logger.Info("get balance cron begin")
		getHuobiBalance()
		logger.Info("get balance cron end")
	})
}

func getHuobiBalance() {
	strategys, err := models.GetAllStrategys()
	if err != nil {
		logger.Error("get strategys err:", err)
		return
	}
	usdt := 0.0
	btc := 0.0
	for _, s := range strategys {
		usdt += huobi.GetCurrencyBalance(s, "usdt")
		btc += huobi.GetCurrencyBalance(s, "btc")
	}
	balanceMap := make(map[string]float64)
	balanceMap["usdt"] = usdt
	balanceMap["btc"] = btc
	b, _ := json.Marshal(balanceMap)
	logger.Info("balace map", string(b))
	if conf.Config.Redis.IsUse {
		redis.GoRedisClient.Set(consts.HUOBI_BALANCE_KEY, string(b), time.Second*consts.REDIS_KEY_EXPIRED_SECONDS)
	} else {
		logger.Error("redis client error...")
	}
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
