// @ticket 02-支付宝支付：沙箱客户端初始化、支付链接生成、异步通知验签
package pkg

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	alipay "github.com/smartwalle/alipay/v3"

	"smart_lockers/internal/config"
)

var (
	aliMu     sync.RWMutex
	aliClient *alipay.Client
)

// InitAlipay 初始化支付宝客户端；sandbox=true 自动使用官方沙箱网关 @ticket 02-支付宝支付
func InitAlipay(cfg config.AlipayConfig) error {
	if cfg.AppID == "" || cfg.PrivateKey == "" {
		return fmt.Errorf("alipay app_id/private_key is empty")
	}
	c, err := alipay.New(cfg.AppID, cfg.PrivateKey, !cfg.Sandbox)
	if err != nil {
		return fmt.Errorf("alipay client init: %w", err)
	}
	if cfg.PublicKey != "" {
		if err = c.LoadAliPayPublicKey(cfg.PublicKey); err != nil {
			return fmt.Errorf("alipay public key load: %w", err)
		}
	}
	aliMu.Lock()
	aliClient = c
	aliMu.Unlock()
	return nil
}

// AlipayReady 客户端是否已就绪 @ticket 02-支付宝支付
func AlipayReady() bool {
	aliMu.RLock()
	defer aliMu.RUnlock()
	return aliClient != nil
}

// CreateAlipayPagePayURL 电脑网站支付：本地签名生成跳转链接（无出站调用）@ticket 02-支付宝支付
func CreateAlipayPagePayURL(cfg config.AlipayConfig, outTradeNo string, amountFen int64, subject string) (string, error) {
	aliMu.RLock()
	c := aliClient
	aliMu.RUnlock()
	if c == nil {
		return "", fmt.Errorf("alipay client not initialized")
	}
	if outTradeNo == "" || amountFen <= 0 {
		return "", fmt.Errorf("invalid out_trade_no or amount")
	}
	p := alipay.TradePagePay{Trade: alipay.Trade{
		Subject:     subject,
		OutTradeNo:  outTradeNo,
		TotalAmount: fmt.Sprintf("%.2f", float64(amountFen)/100.0),
		ProductCode: "FAST_INSTANT_TRADE_PAY",
		NotifyURL:   cfg.NotifyURL,
		ReturnURL:   cfg.ReturnURL,
	}}
	u, err := c.TradePagePay(p)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// VerifyAlipayNotify 校验支付宝异步通知签名（RSA2）@ticket 02-支付宝支付
func VerifyAlipayNotify(values url.Values) error {
	aliMu.RLock()
	c := aliClient
	aliMu.RUnlock()
	if c == nil {
		return fmt.Errorf("alipay client not initialized")
	}
	return c.VerifySign(context.Background(), values)
}
