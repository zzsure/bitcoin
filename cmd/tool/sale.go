package tool

import (
	"encoding/json"
	"errors"

	"github.com/op/go-logging"
	"github.com/urfave/cli"
	"gitlab.azbit.cn/web/bitcoin/cmd"
	"gitlab.azbit.cn/web/bitcoin/conf"
	"gitlab.azbit.cn/web/bitcoin/library/db"
	"gitlab.azbit.cn/web/bitcoin/library/util"
	"gitlab.azbit.cn/web/bitcoin/models"
	"gitlab.azbit.cn/web/bitcoin/modules/huobi"
)

var logger = logging.MustGetLogger("cmd/tool")

var Sale = cli.Command{
	Name:  "sale",
	Usage: "bitcoin sale current orders",
	Flags: []cli.Flag{
		cmd.StringFlag("conf, c", "config.toml", "toml配置文件"),
		cmd.StringFlag("args, a", "", "cmd line args"),
		cmd.IntFlag("strategy", 0, "sale strategy id"),
	},
	Action: runSale,
}

func runSale(c *cli.Context) {
	conf.Init(c.String("conf"), c.String("args"))
	db.Init()
	strategyID := c.Int("strategy")
	err := saleStrategy(strategyID)
	if err != nil {
		logger.Error("sale err:", err)
	}
}

func saleStrategy(strategyID int) error {
	// TODO:清除正在运行的策略
	strategy, err := models.GetStrategyByID(strategyID)
	if err != nil {
		return err
	}
	ol, err := models.GetOrdersByStatus(uint(strategyID), models.OrderStatusSuccess)
	if err != nil {
		return err
	}

	amount := 0.0
	for _, o := range ol {
		if o.Type == models.OrderTypeBuy {
			amount += o.Amount
		} else if o.Type == models.OrderTypeSale {
			amount -= o.Amount
		}
	}

	amount = util.Float64Precision(amount, 4, false)
	if amount <= 0.0 {
		return errors.New("no order have amount")
	}
	rp := 0.0
	r := huobi.GetKLine("btcusdt", "1min", 1)
	if len(r.Data) > 0 {
		rp = (r.Data[0].High + r.Data[0].Low) / 2
	}
	o, err := huobi.HuobiPlaceOrder(strategy, "btcusdt", models.OrderTypeSale, amount)
	if err != nil {
		return err
	}
	o.RefrencePrice = rp
	err = o.Save()
	if err != nil {
		return err
	}
	// 结算
	settle(strategyID)
	return err
}

func settle(strategyID int) error {
	strategy, err := models.GetStrategyByID(strategyID)
	if err != nil {
		return err
	}
	ol, err := models.GetOrdersByStatus(uint(strategyID), models.OrderStatusSuccess)
	if err != nil {
		return err
	}
	ids := make([]uint, len(ol))
	reason := "tool_sale"
	capital := 0.0
	income := 0.0
	fee := 0.0
	for idx, o := range ol {
		if o.Type == models.OrderTypeBuy {
			capital += o.Money
		} else if o.Type == models.OrderTypeSale {
			income += o.Money
		}
		fee += o.Fee
		ids[idx] = o.ID
		o.Status = models.OrderStatusSettle
		err := o.Save()
		if err != nil {
			return err
		}
	}
	idsByte, _ := json.Marshal(ids)
	p := &models.Profit{
		StrategyID:  strategy.ID,
		TotalAmount: strategy.TotalAmount,
		Depth:       strategy.Depth,
		FloatRate:   strategy.FloatRate,
		Capital:     capital,
		InCome:      income,
		Fee:         fee,
		Profit:      income - capital,
		Reason:      reason,
		Day:         util.GetTodayDay(),
		Orders:      string(idsByte),
	}
	err = p.Save()
	return err
}
