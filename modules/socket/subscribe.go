package socket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/bitcoin/library/util"
	"gitlab.azbit.cn/web/bitcoin/models"
)

var logger = logging.MustGetLogger("modules/socket")

type PingMsg struct {
	Ping  int       `json:"ping"`
	Ts    int64     `json:"ts"` // 响应生成时间点, 单位毫秒
	Ch    string    `json:"ch"` // 数据所属的Channel, 格式: market.$symbol.kline.$period
	KLine PingKLine `json:"tick"`
}

type PingKLine struct {
	ID     int64   `json:"id"`     // K线ID
	Amount float64 `json:"amount"` // 成交量
	Count  int64   `json:"count"`  // 成交笔数
	Open   float64 `json:"open"`   // 开盘价
	Close  float64 `json:"close"`  // 收盘价, 当K线为最晚的一根时, 时最新成交价
	Low    float64 `json:"low"`    // 最低价
	High   float64 `json:"high"`   // 最高价
	Vol    float64 `json:"vol"`    // 成交额, 即SUM(每一笔成交价 * 该笔的成交数量)
}

type PongMsg struct {
	Pong int `json:"pong"`
}

type SubMsg struct {
	Sub string `json:"sub"`
	Id  string `json:"id"`
}

func Init() {
	c, _, err := websocket.DefaultDialer.Dial("wss://api.huobi.pro/ws", nil)
	if err != nil {
		logger.Error("dial:", err)
	} else {
		logger.Info("connect huobi success")
		id := util.GenServerUUID()
		msg := SubMsg{
			Sub: "market.btcusdt.kline.1min",
			Id:  id,
		}
		str, _ := json.Marshal(msg)
		err = c.WriteMessage(websocket.TextMessage, []byte(str))
		if err != nil {
			logger.Error("write sub:", err)
			return
		}
	}

	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			logger.Error("read:", err)
			return
		}
		message, err = util.UnzipData(message)
		logger.Info("recv: %s", string(message))
		var data PingMsg
		err = json.Unmarshal(message, &data)
		if err != nil {
			logger.Error("json unmarshal:", err)
			return
		}
		if data.Ping > 0 {
			msg := PongMsg{
				Pong: data.Ping,
			}
			str, _ := json.Marshal(msg)
			err := c.WriteMessage(websocket.TextMessage, []byte(str))
			if err != nil {
				logger.Error("write:", err)
				return
			}
		} else {
			kld := &models.KLineData{
				KID:    data.KLine.ID,
				Amount: data.KLine.Amount,
				Count:  data.KLine.Count,
				Open:   data.KLine.Open,
				Close:  data.KLine.Close,
				Low:    data.KLine.Low,
				High:   data.KLine.High,
				Vol:    data.KLine.Vol,
				Ts:     data.Ts,
				Ch:     data.Ch,
			}
			err := kld.Save()
			if err != nil {
				logger.Error("save kline data to db: ", err)
			}
		}
	}
}
