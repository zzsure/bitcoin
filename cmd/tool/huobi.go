package tool

import (
	"encoding/json"

	"github.com/urfave/cli"
	"gitlab.azbit.cn/web/bitcoin/cmd"
	"gitlab.azbit.cn/web/bitcoin/conf"
	"gitlab.azbit.cn/web/bitcoin/library/db"
	"gitlab.azbit.cn/web/bitcoin/models"
	"gitlab.azbit.cn/web/bitcoin/modules/huobi"
)

const (
	HUOBI_API_UNKOWN = iota
	HUOBI_API_ACCOUNT
	HUOBI_API_BALANCE_USDT
	HUOBI_API_BALANCE_BTC
	HUOBI_API_BUY_BTC
	HUOBI_API_KLINE
	HUOBI_API_TRADE
)

var Huobi = cli.Command{
	Name:  "huobi",
	Usage: "test huobi api",
	Flags: []cli.Flag{
		cmd.StringFlag("conf, c", "config.toml", "toml配置文件"),
		cmd.StringFlag("args, a", "", "cmd line args"),
		cmd.IntFlag("type", 0, "api type"),
	},
	Action: runHuobi,
}

func runHuobi(c *cli.Context) {
	conf.Init(c.String("conf"), c.String("args"))
	db.Init()
	apiType := c.Int("type")
	strategys, err := models.GetAllStrategys()
	if err != nil {
		logger.Error("get strategys err:", err)
		return
	}
	if apiType == HUOBI_API_KLINE {
		kline := huobi.GetKLine("btcusdt", "1min", 1)
		klineByte, _ := json.Marshal(kline)
		logger.Info("kline:", string(klineByte))
		return
	}
	for _, s := range strategys {
		if apiType == HUOBI_API_ACCOUNT {
			account := huobi.GetAccounts(s)
			logger.Info("account:", account)
		} else if apiType == HUOBI_API_BALANCE_USDT {
			balance := huobi.GetCurrencyBalance(s, "usdt")
			logger.Info("balance usdt:", balance)
		} else if apiType == HUOBI_API_BALANCE_BTC {
			balance := huobi.GetCurrencyBalance(s, "btc")
			logger.Info("balance btc:", balance)
		} else if apiType == HUOBI_API_BUY_BTC {
			_, err = huobi.HuobiPlaceOrder(s, "btcusdt", models.OrderTypeBuy, 2)
			logger.Info("buy btc:", err)
		} else if apiType == HUOBI_API_TRADE {
			t := huobi.GetOrders(s, "btcusdt", 100)
			j, _ := json.Marshal(t)
			logger.Info("strategy id:", s.ID, ", btc trade:", string(j))
		}
	}
}
