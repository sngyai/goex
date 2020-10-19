package binance

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)


type BinanceMargin struct {
	Binance
}

func NewMargin(client *http.Client, api_key, secret_key string) *BinanceMargin {
	return NewWithConfig(&APIConfig{
		HttpClient:   client,
		Endpoint:     GLOBAL_API_BASE_URL,
		ApiKey:       api_key,
		ApiSecretKey: secret_key})
}

func NewMarginWithConfig(config *APIConfig) *BinanceMargin {
	if config.Endpoint == "" {
		config.Endpoint = GLOBAL_API_BASE_URL
	}

	bn := &BinanceMargin{Binance:Binance{
		baseUrl:    config.Endpoint,
		apiV1:      config.Endpoint + "/sapi/v1/",
		apiV3:      config.Endpoint + "/sapi/v3/",
		accessKey:  config.ApiKey,
		secretKey:  config.ApiSecretKey,
		httpClient: config.HttpClient}}
	bn.setTimeOffset()
	return bn
}

func (bn *BinanceMargin) Borrow(isIsolated bool, param BorrowParameter) (int64, error) {
	path := bn.apiV1 + "margin/loan"
	params := url.Values{}
	params.Set("symbol", param.CurrencyPair.ToSymbol(""))
	params.Set("asset", param.Currency.String())
	params.Set("isIsolated", strconv.FormatBool(isIsolated))
	params.Set("amount", fmt.Sprint(param.Amount))

	bn.buildParamsSigned(&params)

	resp, err := HttpPostForm2(bn.httpClient, path, params,
		map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return 0, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return 0, err
	}

	orderId := ToInt64(respmap["tranId"])
	if orderId <= 0 {
		return 0, errors.New(string(resp))
	}
	return orderId, nil
}
