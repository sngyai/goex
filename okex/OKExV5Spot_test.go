package okex

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nntaoli-project/goex"
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
		Endpoint:      "https://www.okex.com",
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
