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
// 定义函数 Setup：封装一段可复用的处理逻辑（参数：）

func Setup() *gin.Engine {
	// 执行 gin.SetMode(config.Get().Server.Mode)，完成对应的业务处理

	gin.SetMode(config.Get().Server.Mode)
	// 定义并初始化局部变量 r，取值为 gin.New()

	r := gin.New()
	// 执行 r.Use(middleware.Recovery(), middleware.RequestID(), middleware.Logger(), middleware.CORS())，完成对应的业务处理

	r.Use(middleware.Recovery(), middleware.RequestID(), middleware.Logger(), middleware.CORS())

	// 执行 r.GET("/health", func(c *gin.Context) { pkg.OK(c, map[string]interface{}{"status": "up"}) })，完成对应的业务处理

	r.GET("/health", func(c *gin.Context) { pkg.OK(c, map[string]interface{}{"status": "up"}) })

	// 工单07 数字人 WebSocket：本地演示模式，前端 Vite 代理 /ws → ws://localhost:18080 @ticket 07-数字人Coze

	r.GET("/ws/digital-human", handler.DigitalHumanWS)
	// 工单07 数字人账号密码登录：真实查库校验，前端 Vite 代理 /api → http://localhost:18080 @ticket 07-数字人Coze

	r.POST("/api/v1/digital-human/login", handler.DigitalHumanLogin)

	// 定义并初始化局部变量 v1，取值为 r.Group("/api/v1")

	v1 := r.Group("/api/v1")

	// 鉴权（06）
	// 定义并初始化局部变量 auth，取值为 v1.Group("/auth")

	auth := v1.Group("/auth")
	// 进入该代码块，开始执行其中的逻辑

	{
		// 执行 auth.POST("/login", handler.Login)，完成对应的业务处理

		auth.POST("/login", handler.Login)
		// 修复偏差4：注册手机号+验证码演示登录，前端 LoginView 走真实 HTTP 签发 token @ticket 00-project

		auth.POST("/demo-login", handler.DemoLogin)
		// 执行 auth.POST("/logout", middleware.Auth(), handler.Logout)，完成对应的业务处理

		auth.POST("/logout", middleware.Auth(), handler.Logout)
		// 当前代码块结束

	}

	// 工单03 首页（定位/城市/寄存点/指南，读接口不强制鉴权）
	// 定义并初始化局部变量 home，取值为 v1.Group("/home")

	home := v1.Group("/home")
	// 进入该代码块，开始执行其中的逻辑

	{
		// 执行 home.POST("/locate", handler.Locate)，完成对应的业务处理

		home.POST("/locate", handler.Locate)
		// 执行 home.GET("/cities", handler.SearchCities)，完成对应的业务处理

		home.GET("/cities", handler.SearchCities)
		// 执行 home.GET("/sites", handler.NearbySites)，完成对应的业务处理

		home.GET("/sites", handler.NearbySites)
		// 执行 home.GET("/guide", handler.GetGuide)，完成对应的业务处理

		home.GET("/guide", handler.GetGuide)
		// 当前代码块结束

	}

	// 工单04 物品寄存（写接口鉴权）
	// 定义并初始化局部变量 dep，取值为 v1.Group("/deposit")

	dep := v1.Group("/deposit")
	// 进入该代码块，开始执行其中的逻辑

	{
		// 执行 dep.GET("/cabinets", handler.QueryCabinets)，完成对应的业务处理

		dep.GET("/cabinets", handler.QueryCabinets)
		// 执行 dep.GET("/quote", handler.QuoteFee)，完成对应的业务处理

		dep.GET("/quote", handler.QuoteFee)
		// 执行 dep.POST("/orders", middleware.Auth(), handler.CreateDeposit)，完成对应的业务处理

		dep.POST("/orders", middleware.Auth(), handler.CreateDeposit)
		// 执行 dep.POST("/orders/:order_no/change-cabinet", middleware.Auth(), handler.ChangeCabinet)，完成对应的业务处理

		dep.POST("/orders/:order_no/change-cabinet", middleware.Auth(), handler.ChangeCabinet)
		// 当前代码块结束

	}

	// 工单02 支付宝支付（鉴权 + 回调）
	// 定义并初始化局部变量 pay，取值为 v1.Group("/pay")

	pay := v1.Group("/pay")
	// 进入该代码块，开始执行其中的逻辑

	{
		// 执行 pay.GET("/balance", middleware.Auth(), handler.Balance)，完成对应的业务处理

		pay.GET("/balance", middleware.Auth(), handler.Balance)
		// 执行 pay.GET("/balance/logs", middleware.Auth(), handler.BalanceLogs)，完成对应的业务处理

		pay.GET("/balance/logs", middleware.Auth(), handler.BalanceLogs)
		// 执行 pay.POST("/recharge", middleware.Auth(), handler.CreateRecharge)，完成对应的业务处理

		pay.POST("/recharge", middleware.Auth(), handler.CreateRecharge)
		// 执行 pay.GET("/recharge/records", middleware.Auth(), handler.RechargeRecords)，完成对应的业务处理

		pay.GET("/recharge/records", middleware.Auth(), handler.RechargeRecords)
		// 执行 pay.POST("/alipay/callback", handler.RechargeCallback) // 兼容本地 JSON 回调，完成对应的业务处理

		pay.POST("/alipay/callback", handler.RechargeCallback) // 兼容本地 JSON 回调
		// 执行 pay.POST("/alipay/notify", handler.AlipayNotify)       // 支付宝官方异步通知（form + RSA2 验签），完成对应的业务处理

		pay.POST("/alipay/notify", handler.AlipayNotify) // 支付宝官方异步通知（form + RSA2 验签）
		// 当前代码块结束

	}

	// 工单05 订单查看与取出（全部鉴权）
	// 定义并初始化局部变量 orders，取值为 v1.Group("/orders")

	orders := v1.Group("/orders")
	// 执行 orders.Use(middleware.Auth())，完成对应的业务处理

	orders.Use(middleware.Auth())
	// 进入该代码块，开始执行其中的逻辑

	{
		// 执行 orders.GET("", handler.OrderList)，完成对应的业务处理

		orders.GET("", handler.OrderList)
		// 执行 orders.GET("/summary/ongoing", handler.OngoingSummary)，完成对应的业务处理

		orders.GET("/summary/ongoing", handler.OngoingSummary)
		// 执行 orders.GET("/:order_no", handler.OrderDetail)，完成对应的业务处理

		orders.GET("/:order_no", handler.OrderDetail)
		// 执行 orders.POST("/:order_no/open-auth", handler.CreateOpenAuth)，完成对应的业务处理

		orders.POST("/:order_no/open-auth", handler.CreateOpenAuth)
		// 执行 orders.POST("/:order_no/makeup", handler.MakeupPayment)，完成对应的业务处理

		orders.POST("/:order_no/makeup", handler.MakeupPayment)
		// 执行 orders.POST("/:order_no/complete", handler.CompleteOrder)，完成对应的业务处理

		orders.POST("/:order_no/complete", handler.CompleteOrder)
		// 执行 orders.POST("/open-auth/:token/use", handler.UseOpenAuth)，完成对应的业务处理

		orders.POST("/open-auth/:token/use", handler.UseOpenAuth)
		// 当前代码块结束

	}

	// 工单06 个人中心（鉴权）
	// 定义并初始化局部变量 user，取值为 v1.Group("/user")

	user := v1.Group("/user")
	// 执行 user.Use(middleware.Auth())，完成对应的业务处理

	user.Use(middleware.Auth())
	// 进入该代码块，开始执行其中的逻辑

	{
		// 执行 user.GET("/order-summary", handler.OrderSummary)，完成对应的业务处理

		user.GET("/order-summary", handler.OrderSummary)
		// 执行 user.POST("/avatar", handler.UploadAvatar)，完成对应的业务处理

		user.POST("/avatar", handler.UploadAvatar)
		// 执行 user.GET("/services", handler.MyServices)，完成对应的业务处理

		user.GET("/services", handler.MyServices)
		// 执行 user.GET("/coupons", handler.Coupons)，完成对应的业务处理

		user.GET("/coupons", handler.Coupons)
		// 保留 GET 兼容旧调用，补充 POST 与前端 validateCoupon 的 POST 请求对齐 @ticket 06-个人中心

		user.GET("/coupons/:code/validate", handler.ValidateCoupon)
		// 修复偏差1：用户端前端 validateCoupon 使用 POST，补充同路径 POST 路由（与 GET 并存）@ticket 06-个人中心

		user.POST("/coupons/:code/validate", handler.ValidateCoupon)
		// 当前代码块结束

	}

	// 工单01 消息通知（鉴权）
	// 定义并初始化局部变量 notify，取值为 v1.Group("/notify")

	notify := v1.Group("/notify")
	// 执行 notify.Use(middleware.Auth())，完成对应的业务处理

	notify.Use(middleware.Auth())
	// 进入该代码块，开始执行其中的逻辑

	{
		// 执行 notify.POST("/send", handler.Send)，完成对应的业务处理

		notify.POST("/send", handler.Send)
		// 执行 notify.POST("/scene", handler.SendScene)，完成对应的业务处理

		notify.POST("/scene", handler.SendScene)
		// 执行 notify.GET("/scenes", handler.Scenes)，完成对应的业务处理

		notify.GET("/scenes", handler.Scenes)
		// 执行 notify.GET("/records", handler.Records)，完成对应的业务处理

		notify.GET("/records", handler.Records)
		// 执行 notify.GET("/monitor", handler.Monitor)，完成对应的业务处理

		notify.GET("/monitor", handler.Monitor)
		// 当前代码块结束

	}

	// 工单08 性能分析（鉴权内网限制）
	// 定义并初始化局部变量 ops，取值为 v1.Group("/ops")

	ops := v1.Group("/ops")
	// 执行 ops.Use(middleware.Auth())，完成对应的业务处理

	ops.Use(middleware.Auth())
	// 进入该代码块，开始执行其中的逻辑

	{
		// 执行 ops.GET("/pprof/", handler.PprofIndex)，完成对应的业务处理

		ops.GET("/pprof/", handler.PprofIndex)
		// 执行 ops.GET("/pprof/profile", handler.PprofProfile)，完成对应的业务处理

		ops.GET("/pprof/profile", handler.PprofProfile)
		// 执行 ops.GET("/pprof/heap", handler.PprofHeap)，完成对应的业务处理

		ops.GET("/pprof/heap", handler.PprofHeap)
		// 执行 ops.GET("/pprof/goroutine", handler.PprofGoroutine)，完成对应的业务处理

		ops.GET("/pprof/goroutine", handler.PprofGoroutine)
		// 执行 ops.GET("/pprof/block", handler.PprofBlock)，完成对应的业务处理

		ops.GET("/pprof/block", handler.PprofBlock)
		// 执行 ops.GET("/trace", handler.Trace)，完成对应的业务处理

		ops.GET("/trace", handler.Trace)
		// 执行 ops.GET("/expvar", handler.Expvar)，完成对应的业务处理

		ops.GET("/expvar", handler.Expvar)
		// 当前代码块结束

	}

	// IoT 设备接入（工单 IoT-MQTT）：设备侧凭证/Hook 接口免 JWT（凭 SN+Secret/Token 校验），管理接口走 JWT @ticket IoT-MQTT设备接入
	// 定义并初始化局部变量 dev，取值为 v1.Group("/devices")

	dev := v1.Group("/devices")
	// 进入该代码块，开始执行其中的逻辑

	{
		// 设备请求 MQTT 凭证（方案 2.2 认证接口）
		// 执行 dev.POST("/auth", handler.DeviceAuth)，完成对应的业务处理

		dev.POST("/auth", handler.DeviceAuth)
		// EMQX HTTP 认证 Hook（Token 校验）
		// 执行 dev.POST("/emqx/auth", handler.EMQXAuthHook)，完成对应的业务处理

		dev.POST("/emqx/auth", handler.EMQXAuthHook)
		// EMQX 上下线 Webhook
		// 执行 dev.POST("/emqx/webhook", handler.EMQXWebhook)，完成对应的业务处理

		dev.POST("/emqx/webhook", handler.EMQXWebhook)
	// 当前代码块结束

	}
	// IoT 设备管理查询/指令/OTA（JWT 鉴权）@ticket IoT-MQTT设备接入
	// 定义并初始化局部变量 devAuth，取值为 v1.Group("/devices")

	devAuth := v1.Group("/devices")
	// 执行 devAuth.Use(middleware.Auth())，完成对应的业务处理

	devAuth.Use(middleware.Auth())
	// 进入该代码块，开始执行其中的逻辑

	{
		// 设备状态列表
		// 执行 devAuth.GET("/status", handler.DeviceStatusH)，完成对应的业务处理

		devAuth.GET("/status", handler.DeviceStatusH)
		// 下发设备指令（开门/关门/重启/配置）
		// 执行 devAuth.POST("/command", handler.DeviceCommandH)，完成对应的业务处理

		devAuth.POST("/command", handler.DeviceCommandH)
		// 指令记录查询
		// 执行 devAuth.GET("/commands", handler.DeviceCommandsH)，完成对应的业务处理

		devAuth.GET("/commands", handler.DeviceCommandsH)
		// 设备上行事件查询
		// 执行 devAuth.GET("/events", handler.DeviceEventsH)，完成对应的业务处理

		devAuth.GET("/events", handler.DeviceEventsH)
		// 创建 OTA 升级任务
		// 执行 devAuth.POST("/ota", handler.DeviceOtaCreateH)，完成对应的业务处理

		devAuth.POST("/ota", handler.DeviceOtaCreateH)
		// OTA 任务列表
		// 执行 devAuth.GET("/ota/list", handler.DeviceOtaListH)，完成对应的业务处理

		devAuth.GET("/ota/list", handler.DeviceOtaListH)
	// 当前代码块结束

	}

	// 返回 r 给调用方，结束当前函数

	return r
	// 当前代码块结束

}
