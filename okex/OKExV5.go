package okex

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	. "github.com/nntaoli-project/goex"
)

const v5RestBaseUrl = "https://www.okex.com"
const v5WsBaseUrl = "wss://ws.okex.com:8443/ws/v5"

type OKExV5 struct {
	config        *APIConfig
	customCIDFunc func() string
}

func NewOKExV5(config *APIConfig) *OKExV5 {
	if config.Endpoint == "" {
		config.Endpoint = v5RestBaseUrl
	}
	okex := &OKExV5{config: config}
	return okex
}

func (ok *OKExV5) GetExchangeName() string {
	return OKEX
}

func (ok *OKExV5) UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", 32)
}

func (ok *OKExV5) SetCustomCID(f func() string) {
	ok.customCIDFunc = f
}

//获取所有产品行情信息
//产品类型instType
// SPOT：币币
// SWAP：永续合约
// FUTURES：交割合约
// OPTION：期权
// func (ok *OKExV5) GetTickersV5(instType, uly string) ([]Ticker, error) {
// 	urlPath := fmt.Sprintf("/api/v5/market/tickers?instType=%s", instType)
// 	if instType == "SWAP" || instType == "FUTURES" || instType == "OPTION" {
// 		urlPath = fmt.Sprintf("%s&uly=%s", urlPath, uly)
// 	}
// 	var response spotTickerResponse
// 	err := ok.OKEx.DoRequest("GET", urlPath, "", &response)
// 	if err != nil {
// 		return nil, err
// 	}

// 	date, _ := time.Parse(time.RFC3339, response.Timestamp)
// 	return &Ticker{
// 		Pair: currency,
// 		Last: response.Last,
// 		High: response.High24h,
// 		Low:  response.Low24h,
// 		Sell: response.BestAsk,
// 		Buy:  response.BestBid,
// 		Vol:  response.BaseVolume24h,
// 		Date: uint64(time.Duration(date.UnixNano() / int64(time.Millisecond)))}, nil

// }

type TickerV5 struct {
	InstId    string  `json:"instId"`
	Last      float64 `json:"last,string"`
	BuyPrice  float64 `json:"bidPx,string"`
	BuySize   float64 `json:"bidSz,string"`
	SellPrice float64 `json:"askPx,string"`
	SellSize  float64 `json:"askSz,string"`
	Open      float64 `json:"open24h,string"`
	High      float64 `json:"high24h,string"`
	Low       float64 `json:"low24h,string"`
	Vol       float64 `json:"volCcy24h,string"`
	VolQuote  float64 `json:"vol24h,string"`
	Timestamp uint64  `json:"ts,string"` // 单位:ms
}

func (ok *OKExV5) GetTickerV5(instId string) (*TickerV5, error) {
	urlPath := fmt.Sprintf("%s/api/v5/market/ticker?instId=%s", ok.config.Endpoint, instId)
	type TickerV5Response struct {
		Code int        `json:"code,string"`
		Msg  string     `json:"msg"`
		Data []TickerV5 `json:"data"`
	}
	var response TickerV5Response
	err := HttpGet4(ok.config.HttpClient, urlPath, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetTickerV5 error:%s", response.Msg)
	}
	return &response.Data[0], nil
}

type DepthV5 struct {
	Asks      [][]string `json:"asks,string"`
	Bids      [][]string `json:"bids,string"`
	Timestamp uint64     `json:"ts,string"` // 单位:ms
}

func (ok *OKExV5) GetDepthV5(instId string, size int) (*DepthV5, error) {

	urlPath := fmt.Sprintf("%s/api/v5/market/books?instId=%s", ok.config.Endpoint, instId)
	if size > 0 {
		urlPath = fmt.Sprintf("%s&sz=%d", urlPath, size)
	}
	type DepthV5Response struct {
		Code int       `json:"code,string"`
		Msg  string    `json:"msg"`
		Data []DepthV5 `json:"data"`
	}
	var response DepthV5Response
	err := HttpGet4(ok.config.HttpClient, urlPath, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetDepthV5 error:%s", response.Msg)
	}
	return &response.Data[0], nil
}

func (ok *OKExV5) GetKlineRecords(instId string, size int) (*DepthV5, error) {

	urlPath := fmt.Sprintf("%s/api/v5/market/candles?instId=%s", ok.config.Endpoint, instId)
	if size > 0 {
		urlPath = fmt.Sprintf("%s&sz=%d", urlPath, size)
	}
	type DepthV5Response struct {
		Code int       `json:"code,string"`
		Msg  string    `json:"msg"`
		Data []DepthV5 `json:"data"`
	}
	var response DepthV5Response
	err := HttpGet4(ok.config.HttpClient, urlPath, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetDepthV5 error:%s", response.Msg)
	}
	return &response.Data[0], nil
}
