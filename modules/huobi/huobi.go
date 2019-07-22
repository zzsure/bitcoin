package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/op/go-logging"
	"gitlab.azbit.cn/web/bitcoin/conf"
	"gitlab.azbit.cn/web/bitcoin/library/util"
	"gitlab.azbit.cn/web/bitcoin/library/util/huobi"
	"gitlab.azbit.cn/web/bitcoin/models"
)

var logger = logging.MustGetLogger("modules/huobi")

// 批量操作的API下个版本再封装

//------------------------------------------------------------------------------------------
// 交易API

// 获取K线数据
// strSymbol: 交易对, btcusdt, bccbtc......
// strPeriod: K线类型, 1min, 5min, 15min......
// nSize: 获取数量, [1-2000]
// return: KLineReturn 对象
func GetKLine(strSymbol, strPeriod string, nSize int) models.KLineReturn {
	kLineReturn := models.KLineReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["period"] = strPeriod
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/kline"
	strUrl := conf.MARKET_URL + strRequestUrl

	jsonKLineReturn := huobi_util.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonKLineReturn), &kLineReturn)

	return kLineReturn
}

// 获取聚合行情
// strSymbol: 交易对, btcusdt, bccbtc......
// return: TickReturn对象
func GetTicker(strSymbol string) models.TickerReturn {
	tickerReturn := models.TickerReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/detail/merged"
	strUrl := conf.MARKET_URL + strRequestUrl

	jsonTickReturn := huobi_util.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonTickReturn), &tickerReturn)

	return tickerReturn
}

// 获取交易深度信息
// strSymbol: 交易对, btcusdt, bccbtc......
// strType: Depth类型, step0、step1......stpe5 (合并深度0-5, 0时不合并)
// return: MarketDepthReturn对象
func GetMarketDepth(strSymbol, strType string) models.MarketDepthReturn {
	marketDepthReturn := models.MarketDepthReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["type"] = strType

	strRequestUrl := "/market/depth"
	strUrl := conf.MARKET_URL + strRequestUrl

	jsonMarketDepthReturn := huobi_util.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonMarketDepthReturn), &marketDepthReturn)

	return marketDepthReturn
}

// 获取交易细节信息
// strSymbol: 交易对, btcusdt, bccbtc......
// return: TradeDetailReturn对象
func GetTradeDetail(strSymbol string) models.TradeDetailReturn {
	tradeDetailReturn := models.TradeDetailReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/trade"
	strUrl := conf.MARKET_URL + strRequestUrl

	jsonTradeDetailReturn := huobi_util.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonTradeDetailReturn), &tradeDetailReturn)

	return tradeDetailReturn
}

// 批量获取最近的交易记录
// strSymbol: 交易对, btcusdt, bccbtc......
// nSize: 获取交易记录的数量, 范围1-2000
// return: TradeReturn对象
func GetTrade(strSymbol string, nSize int) models.TradeReturn {
	tradeReturn := models.TradeReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/trade"
	strUrl := conf.MARKET_URL + strRequestUrl

	jsonTradeReturn := huobi_util.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonTradeReturn), &tradeReturn)

	return tradeReturn
}

// 获取Market Detail 24小时成交量数据
// strSymbol: 交易对, btcusdt, bccbtc......
// return: MarketDetailReturn对象
func GetMarketDetail(strSymbol string) models.MarketDetailReturn {
	marketDetailReturn := models.MarketDetailReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/detail"
	strUrl := conf.MARKET_URL + strRequestUrl

	jsonMarketDetailReturn := huobi_util.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonMarketDetailReturn), &marketDetailReturn)

	return marketDetailReturn
}

//------------------------------------------------------------------------------------------
// 公共API

// 查询系统支持的所有交易及精度
// return: SymbolsReturn对象
func GetSymbols() models.SymbolsReturn {
	symbolsReturn := models.SymbolsReturn{}

	strRequestUrl := "/v1/common/symbols"
	strUrl := conf.TRADE_URL + strRequestUrl

	jsonSymbolsReturn := huobi_util.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonSymbolsReturn), &symbolsReturn)

	return symbolsReturn
}

// 查询系统支持的所有币种
// return: CurrencysReturn对象
func GetCurrencys() models.CurrencysReturn {
	currencysReturn := models.CurrencysReturn{}

	strRequestUrl := "/v1/common/currencys"
	strUrl := conf.TRADE_URL + strRequestUrl

	jsonCurrencysReturn := huobi_util.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonCurrencysReturn), &currencysReturn)

	return currencysReturn
}

// 查询系统当前时间戳
// return: TimestampReturn对象
func GetTimestamp() models.TimestampReturn {
	timestampReturn := models.TimestampReturn{}

	strRequest := "/v1/common/timestamp"
	strUrl := conf.TRADE_URL + strRequest

	jsonTimestampReturn := huobi_util.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonTimestampReturn), &timestampReturn)

	return timestampReturn
}

//------------------------------------------------------------------------------------------
// 用户资产API

// 查询当前用户的所有账户, 根据包含的私钥查询
// return: AccountsReturn对象
func GetAccounts(strategy models.Strategy) models.AccountsReturn {
	accountsReturn := models.AccountsReturn{}

	strRequest := "/v1/account/accounts"
	jsonAccountsReturn := huobi_util.ApiKeyGet(strategy, make(map[string]string), strRequest)
	logger.Info("get accounts: ", jsonAccountsReturn)
	json.Unmarshal([]byte(jsonAccountsReturn), &accountsReturn)

	return accountsReturn
}

// 根据账户ID查询账户余额
// nAccountID: 账户ID, 不知道的话可以通过GetAccounts()获取, 可以只现货账户, C2C账户, 期货账户
// return: BalanceReturn对象
func GetAccountBalance(strategy models.Strategy) models.BalanceReturn {
	balanceReturn := models.BalanceReturn{}

	strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", strategy.AccountID)
	jsonBanlanceReturn := huobi_util.ApiKeyGet(strategy, make(map[string]string), strRequest)
	json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn)

	return balanceReturn
}

func GetCurrencyBalance(strategy models.Strategy, currency string) float64 {
	br := GetAccountBalance(strategy)
	for _, sa := range br.Data.List {
		if sa.Currency == currency && sa.Type == "trade" {
			return util.StringToFloat64(sa.Balance)
		}
	}
	return 0.0
}

//------------------------------------------------------------------------------------------
// 交易API

// 下单
// placeRequestParams: 下单信息
// return: PlaceReturn对象
func Place(strategy models.Strategy, placeRequestParams models.PlaceRequestParams) models.PlaceReturn {
	placeReturn := models.PlaceReturn{}

	mapParams := make(map[string]string)
	mapParams["account-id"] = placeRequestParams.AccountID
	mapParams["amount"] = placeRequestParams.Amount
	if 0 < len(placeRequestParams.Price) {
		mapParams["price"] = placeRequestParams.Price
	}
	if 0 < len(placeRequestParams.Source) {
		mapParams["source"] = placeRequestParams.Source
	}
	mapParams["symbol"] = placeRequestParams.Symbol
	mapParams["type"] = placeRequestParams.Type

	strRequest := "/v1/order/orders/place"
	jsonPlaceReturn := huobi_util.ApiKeyPost(strategy, mapParams, strRequest)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	return placeReturn
}

// 火币下单，写入订单表
func HuobiPlaceOrder(strategy models.Strategy, symbol string, orderType int, amount float64) (*models.Order, error) {
	huobiOrderType := ""
	// 更正浮点精度
	if models.OrderTypeBuy == orderType {
		huobiOrderType = "buy-market"
		amount = util.Float64Precision(amount, 8, false)
	} else if models.OrderTypeSale == orderType {
		huobiOrderType = "sell-market"
		amount = util.Float64Precision(amount, 4, false)
	}

	// 火币下单
	var placeParams models.PlaceRequestParams
	placeParams.AccountID = strategy.AccountID
	placeParams.Amount = util.Float64ToString(amount)
	//placeParams.Price = util.Float64ToString(price)
	placeParams.Source = "api"
	placeParams.Symbol = symbol
	placeParams.Type = huobiOrderType
	logger.Info("Place order with: ", placeParams)
	placeReturn := Place(strategy, placeParams)
	externalID := ""
	if placeReturn.Status == "ok" {
		logger.Info("Place return: ", placeReturn.Data)
		externalID = placeReturn.Data
	} else {
		return nil, errors.New(placeReturn.ErrMsg)
	}
	// 构造订单结构
	o := &models.Order{
		StrategyID: strategy.ID,
		Amount:     amount,
		Type:       orderType,
		Status:     models.OrderStatusBuy, // 模拟下单即买入
		ExternalID: "",
	}
	err := o.Save()
	if err != nil {
		return nil, err
	}

	// 查询订单详情，TODO:好的检查策略
	retryTimes := 0
	orderDetail := PlaceDetail(strategy, externalID)
	for orderDetail.Data.State != "filled" && retryTimes < 11 {
		logger.Info("wait order not filled:", orderDetail)
		time.Sleep(time.Duration(5) * time.Second)
		orderDetail = PlaceDetail(strategy, externalID)
		logger.Info("get order detail success:", orderDetail)
		retryTimes++
	}
	//o.Price = util.StringToFloat64(orderDetail.Data.Price)
	fieldAmount := util.StringToFloat64(orderDetail.Data.FieldAmount)
	fieldFees := util.StringToFloat64(orderDetail.Data.FieldFees)
	fieldCashAmount := util.StringToFloat64(orderDetail.Data.FieldCashAmount)
	if models.OrderTypeBuy == orderType {
		o.Money = fieldCashAmount
		// 已成交数量-已成交手续费
		o.Amount = fieldAmount - fieldFees
		// 已成交总金额/已成交数量
		o.Price = fieldCashAmount / fieldAmount
		o.Fee = fieldCashAmount * fieldFees / fieldAmount
	} else if models.OrderTypeSale == orderType {
		o.Money = fieldCashAmount - fieldFees
		o.Amount = fieldAmount
		o.Price = fieldCashAmount / fieldAmount
		o.Fee = fieldFees
	}
	o.Ts = orderDetail.Data.CreatedAt
	o.ExternalID = externalID
	if orderDetail.Data.State == "filled" {
		o.Status = models.OrderStatusSuccess
	}
	err = o.Save()
	if err != nil {
		return nil, err
	}
	return o, nil
}

// 申请撤销一个订单请求
// strOrderID: 订单ID
// return: PlaceReturn对象
func SubmitCancel(strategy models.Strategy, strOrderID string) models.PlaceReturn {
	placeReturn := models.PlaceReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", strOrderID)
	jsonPlaceReturn := huobi_util.ApiKeyPost(strategy, make(map[string]string), strRequest)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	return placeReturn
}

// 查询订单详情
func PlaceDetail(strategy models.Strategy, strOrderID string) models.PlaceDetailReturn {
	placeDetailReturn := models.PlaceDetailReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s", strOrderID)
	jsonPlaceReturn := huobi_util.ApiKeyGet(strategy, make(map[string]string), strRequest)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeDetailReturn)

	return placeDetailReturn
}
