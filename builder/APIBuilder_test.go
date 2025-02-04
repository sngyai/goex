package builder

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sngyai/goex"
	"github.com/sngyai/goex/internal/logger"
)

var builder = NewAPIBuilder()

func init() {
	logger.SetLevel(logger.INFO)
}

func TestAPIBuilder_Build(t *testing.T) {
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.OKCOIN_COM).GetExchangeName(), goex.OKCOIN_COM)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.HUOBI_PRO).GetExchangeName(), goex.HUOBI_PRO)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.ZB).GetExchangeName(), goex.ZB)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.BIGONE).GetExchangeName(), goex.BIGONE)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.OKEX).GetExchangeName(), goex.OKEX)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.POLONIEX).GetExchangeName(), goex.POLONIEX)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.KRAKEN).GetExchangeName(), goex.KRAKEN)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(goex.FCOIN_MARGIN).GetExchangeName(), goex.FCOIN_MARGIN)
	assert.Equal(t, builder.APIKey("").APISecretkey("").BuildFuture(goex.HBDM).GetExchangeName(), goex.HBDM)
}

func TestAPIBuilder_BuildSpotWs(t *testing.T) {
	//os.Setenv("HTTPS_PROXY" , "socks5://127.0.0.1:1080")
	wsApi, _ := builder.BuildSpotWs(goex.OKEX_V3)
	wsApi.DepthCallback(func(depth *goex.Depth) {
		log.Println(depth)
	})
	wsApi.SubscribeDepth(goex.BTC_USDT)
	time.Sleep(time.Minute)
}

func TestAPIBuilder_BuildFuturesWs(t *testing.T) {
	//os.Setenv("HTTPS_PROXY" , "socks5://127.0.0.1:1080")
	wsApi, _ := builder.BuildFuturesWs(goex.OKEX_V3)
	wsApi.DepthCallback(func(depth *goex.Depth) {
		log.Println(depth)
	})
	wsApi.SubscribeDepth(goex.BTC_USD, goex.QUARTER_CONTRACT)
	time.Sleep(time.Minute)
}

func TestAPIBuilder_GetTicker(t *testing.T) {
	okx := builder.APIKey("").APISecretkey("").Build(goex.OKEX)
	ticker, err := okx.GetTicker(goex.BTC_USDT)
	if err != nil {
		t.Error(err)
	}
	t.Log(ticker)
}
