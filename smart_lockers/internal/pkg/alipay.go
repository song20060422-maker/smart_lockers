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

// 开始定义一组包级/局部变量的 var 声明块

var (
	// 执行 aliMu     sync.RWMutex，完成对应的业务处理

	aliMu     sync.RWMutex
	// 执行 aliClient *alipay.Client，完成对应的业务处理

	aliClient *alipay.Client
// 结束参数/表达式列表

)

// InitAlipay 初始化支付宝客户端；sandbox=true 自动使用官方沙箱网关 @ticket 02-支付宝支付
// 定义函数 InitAlipay：封装一段可复用的处理逻辑（参数：cfg config.AlipayConfig）

func InitAlipay(cfg config.AlipayConfig) error {
	// 条件判断：cfg.AppID == "" || cfg.PrivateKey == "" 时执行下面逻辑

	if cfg.AppID == "" || cfg.PrivateKey == "" {
		// 返回 fmt.Errorf("alipay app_id/private_key is empty") 给调用方，结束当前函数

		return fmt.Errorf("alipay app_id/private_key is empty")
	// 当前代码块结束

	}
	// 定义并初始化局部变量 c, err，取值为 alipay.New(cfg.AppID, cfg.PrivateKey, !cfg.Sandbox)

	c, err := alipay.New(cfg.AppID, cfg.PrivateKey, !cfg.Sandbox)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 fmt.Errorf("alipay client init: %w", err) 给调用方，结束当前函数

		return fmt.Errorf("alipay client init: %w", err)
	// 当前代码块结束

	}
	// 条件判断：cfg.PublicKey != "" 时执行下面逻辑

	if cfg.PublicKey != "" {
		// 条件判断：err = c.LoadAliPayPublicKey(cfg.PublicKey); err != nil 时执行下面逻辑

		if err = c.LoadAliPayPublicKey(cfg.PublicKey); err != nil {
			// 返回 fmt.Errorf("alipay public key load: %w", err) 给调用方，结束当前函数

			return fmt.Errorf("alipay public key load: %w", err)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 执行 aliMu.Lock()，完成对应的业务处理

	aliMu.Lock()
	// 为 aliClient 赋值：c

	aliClient = c
	// 执行 aliMu.Unlock()，完成对应的业务处理

	aliMu.Unlock()
	// 返回 nil 给调用方，结束当前函数

	return nil
// 当前代码块结束

}

// AlipayReady 客户端是否已就绪 @ticket 02-支付宝支付
// 定义函数 AlipayReady：封装一段可复用的处理逻辑（参数：）

func AlipayReady() bool {
	// 执行 aliMu.RLock()，完成对应的业务处理

	aliMu.RLock()
	// 延迟注册 aliMu.RUnlock()，在当前函数返回前自动执行

	defer aliMu.RUnlock()
	// 返回 aliClient != nil 给调用方，结束当前函数

	return aliClient != nil
// 当前代码块结束

}

// CreateAlipayPagePayURL 电脑网站支付：本地签名生成跳转链接（无出站调用）@ticket 02-支付宝支付
// 定义函数 CreateAlipayPagePayURL：封装一段可复用的处理逻辑（参数：cfg config.AlipayConfig, outTradeNo string, amountFen int64, subject string）

func CreateAlipayPagePayURL(cfg config.AlipayConfig, outTradeNo string, amountFen int64, subject string) (string, error) {
	// 执行 aliMu.RLock()，完成对应的业务处理

	aliMu.RLock()
	// 定义并初始化局部变量 c，取值为 aliClient

	c := aliClient
	// 执行 aliMu.RUnlock()，完成对应的业务处理

	aliMu.RUnlock()
	// 条件判断：c == nil 时执行下面逻辑

	if c == nil {
		// 返回 "", fmt.Errorf("alipay client not initialized") 给调用方，结束当前函数

		return "", fmt.Errorf("alipay client not initialized")
	// 当前代码块结束

	}
	// 条件判断：outTradeNo == "" || amountFen <= 0 时执行下面逻辑

	if outTradeNo == "" || amountFen <= 0 {
		// 返回 "", fmt.Errorf("invalid out_trade_no or amount") 给调用方，结束当前函数

		return "", fmt.Errorf("invalid out_trade_no or amount")
	// 当前代码块结束

	}
	// 定义并初始化局部变量 p，取值为 alipay.TradePagePay{Trade: alipay.Trade{

	p := alipay.TradePagePay{Trade: alipay.Trade{
		// 执行 Subject:     subject,，完成对应的业务处理

		Subject:     subject,
		// 执行 OutTradeNo:  outTradeNo,，完成对应的业务处理

		OutTradeNo:  outTradeNo,
		// 执行 TotalAmount: fmt.Sprintf("%.2f", float64(amountFen)/100.0),，完成对应的业务处理

		TotalAmount: fmt.Sprintf("%.2f", float64(amountFen)/100.0),
		// 执行 ProductCode: "FAST_INSTANT_TRADE_PAY",，完成对应的业务处理

		ProductCode: "FAST_INSTANT_TRADE_PAY",
		// 执行 NotifyURL:   cfg.NotifyURL,，完成对应的业务处理

		NotifyURL:   cfg.NotifyURL,
		// 执行 ReturnURL:   cfg.ReturnURL,，完成对应的业务处理

		ReturnURL:   cfg.ReturnURL,
	// 执行 }}，完成对应的业务处理

	}}
	// 定义并初始化局部变量 u, err，取值为 c.TradePagePay(p)

	u, err := c.TradePagePay(p)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 "", err 给调用方，结束当前函数

		return "", err
	// 当前代码块结束

	}
	// 返回 u.String(), nil 给调用方，结束当前函数

	return u.String(), nil
// 当前代码块结束

}

// VerifyAlipayNotify 校验支付宝异步通知签名（RSA2）@ticket 02-支付宝支付
// 定义函数 VerifyAlipayNotify：封装一段可复用的处理逻辑（参数：values url.Values）

func VerifyAlipayNotify(values url.Values) error {
	// 执行 aliMu.RLock()，完成对应的业务处理

	aliMu.RLock()
	// 定义并初始化局部变量 c，取值为 aliClient

	c := aliClient
	// 执行 aliMu.RUnlock()，完成对应的业务处理

	aliMu.RUnlock()
	// 条件判断：c == nil 时执行下面逻辑

	if c == nil {
		// 返回 fmt.Errorf("alipay client not initialized") 给调用方，结束当前函数

		return fmt.Errorf("alipay client not initialized")
	// 当前代码块结束

	}
	// 返回 c.VerifySign(context.Background(), values) 给调用方，结束当前函数

	return c.VerifySign(context.Background(), values)
// 当前代码块结束

}
