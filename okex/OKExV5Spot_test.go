package okex

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/sngyai/goex"
)

func newOKExV5SpotClient() *OKExV5Spot {
	return NewOKExV5Spot(&goex.APIConfig{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return &url.URL{
						Scheme: "socks5",
						Host:   "127.0.0.1:1080"}, nil
				},
			},
		},
		Endpoint:      "https://www.okx.com",
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
	})
}

func TestOKExV5Spot_GetTicker(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetTicker(goex.BTC_USDT))
}

func TestOKExV5Spot_GetDepth(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetDepth(5, goex.BTC_USDT))
}

func TestOKExV5SpotGetKlineRecords(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetKlineRecords(goex.BTC_USDT, goex.KLINE_PERIOD_1MIN, 10))
}

func TestOKExV5Spot_LimitBuy(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.LimitBuy("1", "1.0", goex.XRP_USDT))
	//{"code":"0","data":[{"clOrdId":"0bf60374efe445BC258eddf46df044c3","ordId":"305267682086109184","sCode":"0","sMsg":"","tag":""}],"msg":""}}
}

func TestOKExV5Spot_CancelOrder(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.CancelOrder("305267682086109184", goex.XRP_USDT))
}

func TestOKExV5Spot_GetUnfinishOrders(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetUnfinishOrders(goex.XRP_USDT))
}

func TestOKExV5Spot_GetOneOrder(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetOneOrder("305267682086109184", goex.XRP_USDT))
}

func TestOKExV5Spot_GetAccount(t *testing.T) {
	c := newOKExV5SpotClient()
	t.Log(c.GetAccount())
}

func TestOKExV5Spot_GetExchangeName(t *testing.T) {
	t.Log(newOKExV5SpotClient().GetExchangeName())
}

func TestOKExV5Spot_GetCurrenciesPrecision(t *testing.T) {
	t.Log(newOKExV5SpotClient().GetCurrenciesPrecision())
}
