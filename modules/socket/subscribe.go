package socket

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"bitcoin/conf"
	"bitcoin/library/util"
	"bitcoin/models"
	"bitcoin/modules/strategy"
)

var logger = logging.MustGetLogger("modules/socket")

type PingMsg struct {
	Id     string      `json:"id"`
	Ping   int64       `json:"ping"`
	Rep    string      `json:"rep"`
	Status string      `json:"status"`
	Ts     int64       `json:"ts"` // 响应生成时间点, 单位毫秒
	Ch     string      `json:"ch"` // 数据所属的Channel, 格式: market.$symbol.kline.$period
	KLine  PingKLine   `json:"tick"`
	Data   []PingKLine `json:"data"`
}

var lastPingMsg PingMsg

type PingKLine struct {
	Kid    int64   `json:"id"`     // K线ID
	Amount float64 `json:"amount"` // 成交量
	Count  int64   `json:"count"`  // 成交笔数
	Open   float64 `json:"open"`   // 开盘价
	Close  float64 `json:"close"`  // 收盘价, 当K线为最晚的一根时, 时最新成交价
	Low    float64 `json:"low"`    // 最低价
	High   float64 `json:"high"`   // 最高价
	Vol    float64 `json:"vol"`    // 成交额, 即SUM(每一笔成交价 * 该笔的成交数量)
}

type PongMsg struct {
	Pong int64 `json:"pong"`
}

type SubMsg struct {
	Sub string `json:"sub"`
	Id  string `json:"id"`
}

type ReqMsg struct {
	Req  string `json:"req"`
	Id   string `json:"id"`
	From int64  `json:"from"`
	To   int64  `json:"to"`
}

func sendSubMsg(c *websocket.Conn) error {
	id := util.GenUUID()
	msg := SubMsg{
		Sub: conf.Config.KLineData.Symbol,
		Id:  id,
	}
	subMsg, _ := json.Marshal(msg)
	logger.Info("sub str: ", string(subMsg))
	err := c.WriteMessage(websocket.TextMessage, subMsg)
	return err
}

func sendReqMsg(c *websocket.Conn, id int64) error {
	//id := util.GenUUID()
	to := id + 300*60
	idstr := strconv.FormatInt(to, 10)
	if to > conf.Config.KLineData.To {
		to = conf.Config.KLineData.To
	}
	msg := ReqMsg{
		Req:  conf.Config.KLineData.Symbol,
		Id:   idstr,
		From: id,
		To:   to,
	}
	reqMsg, _ := json.Marshal(msg)
	logger.Info("sub req str: ", string(reqMsg))
	err := c.WriteMessage(websocket.TextMessage, reqMsg)
	return err
}

func writePongMsg(c *websocket.Conn, ts int64) error {
	msg := PongMsg{
		Pong: ts,
	}
	pongMsg, _ := json.Marshal(msg)
	logger.Info("send pont msg: ", string(pongMsg))
	err := c.WriteMessage(websocket.TextMessage, pongMsg)
	return err
}

func saveKLineData(data *PingKLine, ch string, ts int64) {
	kld := &models.KLineData{
		Kid:    data.Kid,
		Amount: data.Amount,
		Count:  data.Count,
		Open:   data.Open,
		Close:  data.Close,
		Low:    data.Low,
		High:   data.High,
		Vol:    data.Vol,
		Ts:     ts,
		Ch:     ch,
	}
	err := kld.Save()
	if err != nil {
		logger.Error("save kline data to db: ", err)
	}
}

func dealMsg(c *websocket.Conn) {
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
			err = writePongMsg(c, data.Ping)
			if err != nil {
				logger.Error("write:", err)
				return
			}
		} else {
			if data.Rep != "" && data.Status == "ok" {
				ch := data.Rep
				var ts int64
				if len(data.Data) > 0 {
					for _, kld := range data.Data {
						ts = kld.Kid
						saveKLineData(&kld, ch, ts)
					}
				} else {
					ts, _ = strconv.ParseInt(data.Id, 10, 64)
				}
				ts = ts + 60
				if ts > conf.Config.KLineData.To {
					logger.Info("get all data done")
					return
				}
				time.Sleep(time.Duration(conf.Config.KLineData.Duration) * time.Second)
				err = sendReqMsg(c, ts)
				if err != nil {
					logger.Error("write req:", err)
				}
			} else if data.Ch != "" {
				if lastPingMsg.Ts != 0 && lastPingMsg.KLine.Kid != data.KLine.Kid {
					saveKLineData(&lastPingMsg.KLine, lastPingMsg.Ch, lastPingMsg.KLine.Kid)
				}
				lastPingMsg = data
				// 应用策略
				kld := &models.KLineData{
					Kid:    data.KLine.Kid,
					Amount: data.KLine.Amount,
					Count:  data.KLine.Count,
					Open:   data.KLine.Open,
					Close:  data.KLine.Close,
					Low:    data.KLine.Low,
					High:   data.KLine.High,
					Vol:    data.KLine.Vol,
					Ts:     data.KLine.Kid,
					Ch:     data.Ch,
				}
				strategy.StrategyDeal(kld)
			} else {
				logger.Error("not sub and not req return data")
				//return
			}
		}
	}
}

func Init() {
	c, _, err := websocket.DefaultDialer.Dial("wss://api.huobi.pro/ws", nil)
	if err != nil {
		logger.Error("dial:", err)
		return
	}
	logger.Info("connect huobi success")

	/*err = sendSubMsg(c)
	if err != nil {
		logger.Error("write sub:", err)
	}*/
	err = sendReqMsg(c, conf.Config.KLineData.From)
	if err != nil {
		logger.Error("write req:", err)
	}

	defer c.Close()

	dealMsg(c)
}
