package models

type PlaceRequestParams struct {
	AccountID string `json:"account-id"` // 账户ID
	Amount    string `json:"amount"`     // 限价表示下单数量, 市价买单时表示买多少钱, 市价卖单时表示卖多少币
	Price     string `json:"price"`      // 下单价格, 市价单不传该参数
	Source    string `json:"source"`     // 订单来源, api: API调用, margin-api: 借贷资产交易
	Symbol    string `json:"symbol"`     // 交易对, btcusdt, bccbtc......
	Type      string `json:"type"`       // 订单类型, buy-market: 市价买, sell-market: 市价卖, buy-limit: 限价买, sell-limit: 限价卖
}

type PlaceReturn struct {
	Status  string `json:"status"`
	Data    string `json:"data"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type PlaceDetail struct {
	ID              int64  `json:"id"`                // 订单ID
	Symbol          string `json:"symbol"`            // 交易对, btcusdt, ethbtc, rcneth ...
	AccountID       int64  `json:"account-id"`        // 账户 ID
	Amount          string `json:"amount"`            // 订单数量
	Price           string `json:"price"`             // 订单价格
	CreatedAt       int64  `json:"created-at"`        // 订单创建时间
	Type            string `json:"type"`              // 订单类型, buy-market：市价买, sell-market：市价卖, buy-limit：限价买, sell-limit：限价卖, buy-ioc：IOC买单, sell-ioc：IOC卖单
	FieldAmount     string `json:"field-amount"`      // 已成交数量
	FieldCashAmount string `json:"field-cash-amount"` // 已成交总金额
	FieldFees       string `json:"field-fees"`        // 已成交手续费（买入为币，卖出为钱）
	FinishAt        int64  `json:"finished-at"`       // 订单变为终结态的时间，不是成交时间，包含“已撤单”状态
	Source          string `json:"source"`            // 订单来源
	State           string `json:"state"`             // 订单状态，submitting , submitted 已提交, partial-filled 部分成交, partial-canceled 部分成交撤销, filled 完全成交, canceled 已撤销
	CanceledAt      int64  `json:"canceled-at"`       // 订单撤销时间
}

type PlaceDetailReturn struct {
	Status string      `json:"status"`
	Data   PlaceDetail `json:"data"`
}

type PlaceOrdersReturn struct {
	Status  string        `json:"status"`
	Data    []PlaceDetail `json:"data"`
	ErrCode string        `json:"err-code"`
	ErrMsg  string        `json:"err-msg"`
}
