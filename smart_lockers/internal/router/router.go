// @ticket 00-project + 各工单：全局路由注册，按业务域分组
package router

import (
	"github.com/gin-gonic/gin"

	"smart_lockers/internal/config"
	"smart_lockers/internal/handler"
	"smart_lockers/internal/middleware"
	"smart_lockers/internal/pkg"
)

// Setup 构建路由引擎 @ticket 00-project
func Setup() *gin.Engine {
	gin.SetMode(config.Get().Server.Mode)
	r := gin.New()
	r.Use(middleware.Recovery(), middleware.RequestID(), middleware.Logger(), middleware.CORS())

	r.GET("/health", func(c *gin.Context) { pkg.OK(c, map[string]interface{}{"status": "up"}) })

	v1 := r.Group("/api/v1")

	// 鉴权（06）
	auth := v1.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/logout", middleware.Auth(), handler.Logout)
	}

	// 工单03 首页（定位/城市/寄存点/指南，读接口不强制鉴权）
	home := v1.Group("/home")
	{
		home.POST("/locate", handler.Locate)
		home.GET("/cities", handler.SearchCities)
		home.GET("/sites", handler.NearbySites)
		home.GET("/guide", handler.GetGuide)
	}

	// 工单04 物品寄存（写接口鉴权）
	dep := v1.Group("/deposit")
	{
		dep.GET("/cabinets", handler.QueryCabinets)
		dep.GET("/quote", handler.QuoteFee)
		dep.POST("/orders", middleware.Auth(), handler.CreateDeposit)
		dep.POST("/orders/:order_no/change-cabinet", middleware.Auth(), handler.ChangeCabinet)
	}

	// 工单02 支付宝支付（鉴权 + 回调）
	pay := v1.Group("/pay")
	{
		pay.GET("/balance", middleware.Auth(), handler.Balance)
		pay.GET("/balance/logs", middleware.Auth(), handler.BalanceLogs)
		pay.POST("/recharge", middleware.Auth(), handler.CreateRecharge)
		pay.GET("/recharge/records", middleware.Auth(), handler.RechargeRecords)
		pay.POST("/alipay/callback", handler.RechargeCallback) // 回调走签名校验，非用户令牌
	}

	// 工单05 订单查看与取出（全部鉴权）
	orders := v1.Group("/orders")
	orders.Use(middleware.Auth())
	{
		orders.GET("", handler.OrderList)
		orders.GET("/summary/ongoing", handler.OngoingSummary)
		orders.GET("/:order_no", handler.OrderDetail)
		orders.POST("/:order_no/open-auth", handler.CreateOpenAuth)
		orders.POST("/:order_no/makeup", handler.MakeupPayment)
		orders.POST("/:order_no/complete", handler.CompleteOrder)
		orders.POST("/open-auth/:token/use", handler.UseOpenAuth)
	}

	// 工单06 个人中心（鉴权）
	user := v1.Group("/user")
	user.Use(middleware.Auth())
	{
		user.GET("/order-summary", handler.OrderSummary)
		user.POST("/avatar", handler.UploadAvatar)
		user.GET("/services", handler.MyServices)
		user.GET("/coupons", handler.Coupons)
		user.GET("/coupons/:code/validate", handler.ValidateCoupon)
	}

	// 工单01 消息通知（鉴权）
	notify := v1.Group("/notify")
	notify.Use(middleware.Auth())
	{
		notify.POST("/send", handler.Send)
		notify.POST("/scene", handler.SendScene)
		notify.GET("/scenes", handler.Scenes)
		notify.GET("/records", handler.Records)
		notify.GET("/monitor", handler.Monitor)
	}

	// 工单08 性能分析（鉴权内网限制）
	ops := v1.Group("/ops")
	ops.Use(middleware.Auth())
	{
		ops.GET("/pprof/", handler.PprofIndex)
		ops.GET("/pprof/profile", handler.PprofProfile)
		ops.GET("/pprof/heap", handler.PprofHeap)
		ops.GET("/pprof/goroutine", handler.PprofGoroutine)
		ops.GET("/pprof/block", handler.PprofBlock)
		ops.GET("/trace", handler.Trace)
		ops.GET("/expvar", handler.Expvar)
	}

	return r
}
