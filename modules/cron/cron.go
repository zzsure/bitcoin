package cron

import (
	"bitcoin/conf"
	"bitcoin/consts"
	"bitcoin/library/redis"
	"encoding/json"
	"time"

	"bitcoin/models"
	"bitcoin/modules/huobi"
	"bitcoin/modules/strategy"
	"github.com/op/go-logging"
	"github.com/robfig/cron"
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
		getHuobiKLine("btcusdt", "market.btcusdt.kline.1min")
		getHuobiKLine("ethusdt", "market.ethusdt.kline.1min")
		logger.Info("get huobi kline cron end")
	})
}

func getBalanceCron() {
	getHuobiBalance()
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

	btcPrice := 0.0
	ethPrice := 0.0
	r := huobi.GetKLine("btcusdt", "1min", 1)
	if len(r.Data) > 0 {
		data := r.Data[0]
		// TODO: price may be 0.0
		btcPrice = data.Open
	}
	e := huobi.GetKLine("ethusdt", "1min", 1)
	if len(e.Data) > 0 {
		data := e.Data[0]
		// TODO: price may be 0.0
		ethPrice = data.Open
	}

	allUsdt := 0.0
	allBtc := 0.0
	allEth := 0.0

	for _, s := range strategys {
		s.UsdtBalance = huobi.GetCurrencyBalance(s, "usdt")
		allUsdt += s.UsdtBalance
		s.BtcBalance = huobi.GetCurrencyBalance(s, "btc")
		allBtc += s.BtcBalance
		s.EthBalance = huobi.GetCurrencyBalance(s, "eth")
		allEth += s.EthBalance
		if btcPrice != 0.0 {
			s.RmbValue = (s.UsdtBalance + s.BtcBalance * btcPrice + s.EthBalance * ethPrice) * 7
			logger.Info("strategy id:", s.ID, ", usdt:", s.UsdtBalance, ", btc:", s.BtcBalance, ", btcPrice:", btcPrice, ", ethPrice:", ethPrice, ", rmb: ", s.RmbValue)
		}
		err := s.Save()
		if err != nil {
			logger.Error("save strategy id: ", s.ID, ", err: ", err)
		}
	}
	balanceMap := make(map[string]float64)
	balanceMap["usdt"] = allUsdt
	balanceMap["btc"] = allBtc
	balanceMap["eth"] = allEth
	b, _ := json.Marshal(balanceMap)
	logger.Info("balace map", string(b))
	if conf.Config.Redis.IsUse {
		redis.GoRedisClient.Set(consts.HUOBI_BALANCE_KEY, string(b), time.Second*consts.REDIS_KEY_EXPIRED_SECONDS)
	} else {
		logger.Error("redis client error...")
	}
}

func getHuobiKLine(symbol, ch string) {
	r := huobi.GetKLine(symbol, "1min", 1)
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
			Ch:     ch,
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
