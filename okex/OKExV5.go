package okex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	. "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
)

const v5RestBaseUrl = "https://www.okex.com"
const v5WsBaseUrl = "wss://ws.okex.com:8443/ws/v5"

// base interface for okex v5
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

func (ok *OKExV5) ExchangeName() string {
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

func (ok *OKExV5) GetKlineRecordsV5(instId, after, before, bar, limit string) ([][]string, error) {

	urlPath := fmt.Sprintf("%s/api/v5/market/candles?instId=%s", ok.config.Endpoint, instId)
	params := url.Values{}
	if after != "" {
		params.Set("after", after)
	}
	if before != "" {
		params.Set("before", before)
	}
	if bar != "" {
		params.Set("bar", bar)
	}
	if limit != "" {
		params.Set("limit", limit)
	}
	if params.Encode() != "" {
		urlPath = fmt.Sprintf("%s&%s", urlPath, params.Encode())
	}

	type CandleResponse struct {
		Code int        `json:"code,string"`
		Msg  string     `json:"msg"`
		Data [][]string `json:"data"`
	}
	var response CandleResponse
	err := HttpGet4(ok.config.HttpClient, urlPath, nil, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("GetKlineRecordsV5 error:%s", response.Msg)
	}
	return response.Data, nil
}

/*
 Get a iso time
  eg: 2018-03-16T18:02:48.284Z
*/
func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

/*
 Get a http request body is a json string and a byte array.
*/
func (ok *OKExV5) BuildRequestBody(params interface{}) (string, *bytes.Reader, error) {
	if params == nil {
		return "", nil, errors.New("illegal parameter")
	}
	data, err := json.Marshal(params)
	if err != nil {
		//log.Println(err)
		return "", nil, errors.New("json convert string error")
	}

	jsonBody := string(data)
	binBody := bytes.NewReader(data)

	return jsonBody, binBody, nil
}

func (ok *OKExV5) doParamSign(httpMethod, uri, requestBody string) (string, string) {
	timestamp := IsoTime()
	preText := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), uri, requestBody)
	//log.Println("preHash", preText)
	sign, _ := GetParamHmacSHA256Base64Sign(ok.config.ApiSecretKey, preText)
	return sign, timestamp
}

func (ok *OKExV5) DoRequest(httpMethod, uri, reqBody string, response interface{}) error {
	url := ok.config.Endpoint + uri
	sign, timestamp := ok.doParamSign(httpMethod, uri, reqBody)
	//logger.Log.Debug("timestamp=", timestamp, ", sign=", sign)
	resp, err := NewHttpRequest(ok.config.HttpClient, httpMethod, url, reqBody, map[string]string{
		CONTENT_TYPE: APPLICATION_JSON_UTF8,
		ACCEPT:       APPLICATION_JSON,
		//COOKIE:               LOCALE + "en_US",
		OK_ACCESS_KEY:        ok.config.ApiKey,
		OK_ACCESS_PASSPHRASE: ok.config.ApiPassphrase,
		OK_ACCESS_SIGN:       sign,
		OK_ACCESS_TIMESTAMP:  fmt.Sprint(timestamp)})
	if err != nil {
		//log.Println(err)
		return err
	} else {
		logger.Log.Debug(string(resp))
		return json.Unmarshal(resp, &response)
	}
}

type CreateOrderParam struct {
	Symbol    string //产品ID
	TradeMode string //交易模式,	保证金模式：isolated：逐仓 ；cross：全仓,	非保证金模式：cash：非保证金
	Side      string // 订单方向 buy：买 sell：卖
	OrderType string //订单类型
	// market：市价单
	// limit：限价单
	// post_only：只做maker单
	// fok：全部成交或立即取消
	// ioc：立即成交并取消剩余

	Size        string //	委托数量
	PosSide     string //持仓方向 在双向持仓模式下必填，且仅可选择 long 或 short
	Price       string //委托价格，仅适用于限价单
	CCY         string // 保证金币种，仅适用于单币种保证金模式下的全仓杠杆订单
	ClientOrdId string //客户自定义订单ID	字母（区分大小写）与数字的组合，可以是纯字母、纯数字且长度要在1-32位之间。
	Tag         string //订单标签	字母（区分大小写）与数字的组合，可以是纯字母、纯数字，且长度在1-8位之间。
	ReduceOnly  bool   //是否只减仓，true 或 false，默认false	仅适用于币币杠杆订单
}

type OrderV5 struct {
	OrdId       string `json:"ordId"`
	ClientOrdId string `json:"clOrdId"` //客户自定义订单ID	字母（区分大小写）与数字的组合，可以是纯字母、纯数字且长度要在1-32位之间。
	Tag         string `json:"tag"`
	SCode       string `json:"sCode"`
	SMsg        string `json:"sMsg"`
}

func (ok *OKExV5) CreateOrder(param *CreateOrderParam) (*OrderV5, error) {

	reqBody := make(map[string]interface{})

	reqBody["instId"] = param.Symbol
	reqBody["tdMode"] = param.TradeMode
	reqBody["side"] = param.Side
	reqBody["ordType"] = param.OrderType
	reqBody["sz"] = param.Size

	if param.CCY != "" {
		reqBody["ccy"] = param.CCY
	}
	if param.ClientOrdId != "" {
		reqBody["clOrdId"] = param.ClientOrdId
	} else {
		if ok.customCIDFunc != nil {
			param.ClientOrdId = ok.customCIDFunc()
		} else {
			param.ClientOrdId = ("0bf60374efe445BC" + strings.Replace(uuid.New().String(), "-", "", 32))[:32]
		}
	}
	if param.Tag != "" {
		reqBody["tag"] = param.Tag
	}
	if param.PosSide != "" {
		reqBody["posSide"] = param.PosSide
	}
	if param.Price != "" {
		reqBody["px"] = param.Price
	}
	if param.ReduceOnly != false {
		reqBody["reduceOnly"] = param.ReduceOnly
	}

	type OrderResponse struct {
		Code int     `json:"code,string"`
		Msg  string  `json:"msg"`
		Data OrderV5 `json:"data"`
	}
	var response OrderResponse

	uri := "/api/v5/trade/order"

	jsonStr, _, _ := ok.BuildRequestBody(param)
	err := ok.DoRequest("POST", uri, jsonStr, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("CreateOrder error:%s", response.Msg)
	}
	return &response.Data, nil
}
