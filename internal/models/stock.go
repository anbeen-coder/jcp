package models

// Stock 股票基本信息
type Stock struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
	Volume        int64   `json:"volume"`
	Amount        float64 `json:"amount"`
	MarketCap     string  `json:"marketCap"`
	Sector        string  `json:"sector"`
	Open          float64 `json:"open"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	PreClose      float64 `json:"preClose"`
}

// KLineData K线数据
type KLineData struct {
	Time   string  `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
	Amount float64 `json:"amount,omitempty"`
	Avg    float64 `json:"avg,omitempty"` // 分时均价线
	// 均线数据
	MA5  float64 `json:"ma5,omitempty"`
	MA10 float64 `json:"ma10,omitempty"`
	MA20 float64 `json:"ma20,omitempty"`
}

// OrderBookItem 盘口单项
type OrderBookItem struct {
	Price   float64 `json:"price"`
	Size    int64   `json:"size"`
	Total   int64   `json:"total"`
	Percent float64 `json:"percent"`
}

// OrderBook 盘口数据
type OrderBook struct {
	Bids []OrderBookItem `json:"bids"`
	Asks []OrderBookItem `json:"asks"`
}

// MarketIndex 大盘指数数据
type MarketIndex struct {
	Code          string  `json:"code"`          // 指数代码，如 sh000001
	Name          string  `json:"name"`          // 指数名称，如 上证指数
	Price         float64 `json:"price"`         // 当前点位
	Change        float64 `json:"change"`        // 涨跌点数
	ChangePercent float64 `json:"changePercent"` // 涨跌幅(%)
	Volume        int64   `json:"volume"`        // 成交量(手)
	Amount        float64 `json:"amount"`        // 成交额(万元)
}
