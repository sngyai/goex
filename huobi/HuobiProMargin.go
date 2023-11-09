package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	. "github.com/sngyai/goex"
)

type HuoBiProMargin struct {
	httpClient *http.Client
	baseUrl    string
	accountId  string
	accessKey  string
	secretKey  string
}

func NewHuoBiProMargin(client *http.Client, apikey, secretkey, accountId string) *HuoBiProMargin {
	hbpro := new(HuoBiProMargin)
	hbpro.baseUrl = "https://api.huobi.pro"
	hbpro.httpClient = client
	hbpro.accessKey = apikey
	hbpro.secretKey = secretkey
	hbpro.accountId = accountId
	return hbpro
}

func (hbpro *HuoBiProMargin) GetFixedAccountInfo(symbol string) (AccountInfo, error) {
	path := "/v1/margin/accounts/balance"
	params := &url.Values{}
	params.Set("symbol", symbol)

	hbpro.buildPostForm("GET", path, params)

	//log.Println(hbpro.BaseUrl + path + "?" + params.Encode())

	respmap, err := HttpGet(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return AccountInfo{}, err
	}

	if respmap["status"].(string) != "ok" {
		return AccountInfo{}, errors.New(respmap["err-code"].(string))
	}
	var info AccountInfo
	return info, nil
}

func (hbpro *HuoBiProMargin) FixedAccountBorrow(param BorrowParameter) (borrowId int64, err error) {
	path := fmt.Sprintf("/v1/margin/orders")
	params := url.Values{}
	params.Set("symbol", strings.ToLower(param.CurrencyPair.ToSymbol("")))
	params.Set("currency", strings.ToLower(param.Currency.String()))
	params.Set("amount", FloatToString(param.Amount, 8))

	hbpro.buildPostForm("POST", path, &params)
	resp, err := HttpPostForm3(hbpro.httpClient, hbpro.baseUrl+path+"?"+params.Encode(), hbpro.toJson(params),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	if err != nil {
		return -1, err
	}

	var respmap map[string]interface{}
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		return -1, err
	}

	if respmap["status"].(string) != "ok" {
		return -1, errors.New(string(resp))
	}

	return ToInt64(respmap["data"]), nil
}

func (hbpro *HuoBiProMargin) buildPostForm(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("AccessKeyId", hbpro.accessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	domain := strings.Replace(hbpro.baseUrl, "https://", "", len(hbpro.baseUrl))
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", reqMethod, domain, path, postForm.Encode())
	sign, _ := GetParamHmacSHA256Base64Sign(hbpro.secretKey, payload)
	postForm.Set("Signature", sign)

	return nil
}

func (hbpro *HuoBiProMargin) toJson(params url.Values) string {
	parammap := make(map[string]string)
	for k, v := range params {
		parammap[k] = v[0]
	}
	jsonData, _ := json.Marshal(parammap)
	return string(jsonData)
}
