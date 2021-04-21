package okex

import (
	"strconv"
	"time"

	. "github.com/nntaoli-project/goex"
)

type OKExV5Spot struct {
	*OKExV5
}

func NewOKExV5Spot(config *APIConfig) *OKExV5Spot {
	if config.Endpoint == "" {
		config.Endpoint = v5RestBaseUrl
	}
	okex := &OKExV5Spot{OKExV5: NewOKExV5(config)}
	return okex
}

// private API
func (ok *OKExV5Spot) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not support")

}
func (ok *OKExV5Spot) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not support")

}
func (ok *OKExV5Spot) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not support")

}
func (ok *OKExV5Spot) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not support")

}
func (ok *OKExV5Spot) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not support")

}
func (ok *OKExV5Spot) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not support")

}
func (ok *OKExV5Spot) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not support")

}
func (ok *OKExV5Spot) GetOrderHistorys(currency CurrencyPair, opt ...OptionalParameter) ([]Order, error) {
	panic("not support")

}
func (ok *OKExV5Spot) GetAccount() (*Account, error) {
	panic("not support")

}

// public API

func (ok *OKExV5Spot) GetTicker(currency CurrencyPair) (*Ticker, error) {
	ticker, err := ok.GetTickerV5(currency.ToSymbol("-"))
	if err != nil {
		return nil, err
	}
	return &Ticker{
		Pair: currency,
		Last: ticker.Last,
		Buy:  ticker.BuyPrice,
		Sell: ticker.SellPrice,
		High: ticker.High,
		Low:  ticker.Low,
		Vol:  ticker.Vol,
		Date: ticker.Timestamp,
	}, nil
}

func (ok *OKExV5Spot) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	d, err := ok.GetDepthV5(currency.ToSymbol("-"), size)
	if err != nil {
		return nil, err
	}

	depth := &Depth{}

	for _, ask := range d.Asks {
		depth.AskList = append(depth.AskList, DepthRecord{Price: ToFloat64(ask[0]), Amount: ToFloat64(ask[1])})
	}
	for _, bid := range d.Bids {
		depth.BidList = append(depth.BidList, DepthRecord{Price: ToFloat64(bid[0]), Amount: ToFloat64(bid[1])})
	}
	depth.UTime = time.Unix(0, int64(d.Timestamp)*1000000)
	return depth, nil
}

func (ok *OKExV5Spot) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error) {
	// [1m/3m/5m/15m/30m/1H/2H/4H/6H/12H/1D/1W/1M/3M/6M/1Y]
	bar := "1D"
	switch period {
	case KLINE_PERIOD_1MIN:
		bar = "1m"
	case KLINE_PERIOD_3MIN:
		bar = "3m"
	case KLINE_PERIOD_5MIN:
		bar = "5m"
	case KLINE_PERIOD_15MIN:
		bar = "15m"
	case KLINE_PERIOD_30MIN:
		bar = "30m"
	case KLINE_PERIOD_1H, KLINE_PERIOD_60MIN:
		bar = "1H"
	case KLINE_PERIOD_2H:
		bar = "2H"
	case KLINE_PERIOD_4H:
		bar = "4H"
	case KLINE_PERIOD_6H:
		bar = "6H"
	case KLINE_PERIOD_12H:
		bar = "12H"
	case KLINE_PERIOD_1DAY:
		bar = "1D"
	case KLINE_PERIOD_1WEEK:
		bar = "1W"
	default:
		bar = "1D"
	}
	after, before, limit := "", "", strconv.Itoa(size)

	for _, opt := range optional {
		for k, v := range opt {
			if k == "after" {
				after = v.(string)
			}
			if k == "before" {
				before = v.(string)
			}
		}
	}
	kl, err := ok.GetKlineRecordsV5(currency.ToSymbol("-"), after, before, bar, limit)
	if err != nil {
		return nil, err
	}

	klines := make([]Kline, 0)

	for _, k := range kl {
		klines = append(klines, Kline{
			Pair:      currency,
			Timestamp: ToInt64(k[0]),
			Open:      ToFloat64(k[1]),
			High:      ToFloat64(k[2]),
			Low:       ToFloat64(k[3]),
			Close:     ToFloat64(k[4]),
			Vol:       ToFloat64(k[5]),
		})
	}

	return klines, nil

}

//非个人，整个交易所的交易记录
func (ok *OKExV5Spot) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not support")
}

func (ok *OKExV5Spot) GetExchangeName() string {
	return ok.ExchangeName() + "_v5_spot"
}
