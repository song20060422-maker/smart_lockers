// @ticket B端-路由注册：按四大功能分组，统一挂 JWT 鉴权 + 角色权限中间件
package router

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/config"
	"smart_lockers_admin/internal/handler"
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

// Setup 构建路由引擎
// 定义函数 Setup：封装一段可复用的处理逻辑（参数：st *store.Store）

func Setup(st *store.Store) *gin.Engine {
	// 执行 gin.SetMode(config.Get().Server.Mode)，完成对应的业务处理

	gin.SetMode(config.Get().Server.Mode)
	// 定义并初始化局部变量 r，取值为 gin.New()

	r := gin.New()
	// 执行 r.Use(middleware.Recovery(), middleware.RequestID(), middleware.Logger(), middleware.CORS())，完成对应的业务处理

	r.Use(middleware.Recovery(), middleware.RequestID(), middleware.Logger(), middleware.CORS())

	// 执行 r.GET("/health", func(c *gin.Context) { pkg.OK(c, map[string]interface{}{"status": "up", "service": "smart_lockers_admin"}) })，完成对应的业务处理

	r.GET("/health", func(c *gin.Context) { pkg.OK(c, map[string]interface{}{"status": "up", "service": "smart_lockers_admin"}) })

	// 定义并初始化局部变量 authH，取值为 &handler.AuthHandler{}

	authH := &handler.AuthHandler{}
	// 定义并初始化局部变量 siteH，取值为 handler.NewSite(serviceSite(st))

	siteH := handler.NewSite(serviceSite(st))
	// 定义并初始化局部变量 devH，取值为 handler.NewDevice(serviceDevice(st))

	devH := handler.NewDevice(serviceDevice(st))
	// 定义并初始化局部变量 recH，取值为 handler.NewReconcile(serviceRecon(st))

	recH := handler.NewReconcile(serviceRecon(st))
	// 定义并初始化局部变量 alarmH，取值为 handler.NewAlarm(serviceAlarm(st))

	alarmH := handler.NewAlarm(serviceAlarm(st))

	// 定义并初始化局部变量 v1，取值为 r.Group("/api/v1")

	v1 := r.Group("/api/v1")

	// 登录：签令牌
	// 定义并初始化局部变量 auth，取值为 v1.Group("/auth")

	auth := v1.Group("/auth")
	// 进入该代码块，开始执行其中的逻辑

	{
		// 执行 auth.POST("/login", authH.Login)，完成对应的业务处理

		auth.POST("/login", authH.Login)
		// 执行 auth.GET("/me", middleware.Auth(), authH.Me)，完成对应的业务处理

		auth.GET("/me", middleware.Auth(), authH.Me)
	// 当前代码块结束

	}

	// 其余接口统一鉴权
	// 定义并初始化局部变量 admin，取值为 v1.Group("/admin")

	admin := v1.Group("/admin")
	// 执行 admin.Use(middleware.Auth())，完成对应的业务处理

	admin.Use(middleware.Auth())
	// 进入该代码块，开始执行其中的逻辑

	{
		// 网点管理：读(3类角色)/写(网点管理员+运营)
		// 定义并初始化局部变量 sites，取值为 admin.Group("/sites")

		sites := admin.Group("/sites")
		// 进入该代码块，开始执行其中的逻辑

		{
			// 执行 sites.GET("", middleware.RequirePerm(middleware.PermSiteRead), siteH.List)，完成对应的业务处理

			sites.GET("", middleware.RequirePerm(middleware.PermSiteRead), siteH.List)
			// 执行 sites.POST("", middleware.RequirePerm(middleware.PermSiteWrite), siteH.Create)，完成对应的业务处理

			sites.POST("", middleware.RequirePerm(middleware.PermSiteWrite), siteH.Create)
			// 执行 sites.POST("/:id/status", middleware.RequirePerm(middleware.PermSiteWrite), siteH.UpdateStatus)，完成对应的业务处理

			sites.POST("/:id/status", middleware.RequirePerm(middleware.PermSiteWrite), siteH.UpdateStatus)
		// 当前代码块结束

		}

		// 设备监控：读/控制
		// 定义并初始化局部变量 dev，取值为 admin.Group("/devices")

		dev := admin.Group("/devices")
		// 进入该代码块，开始执行其中的逻辑

		{
			// 执行 dev.GET("", middleware.RequirePerm(middleware.PermDeviceRead), devH.List)，完成对应的业务处理

			dev.GET("", middleware.RequirePerm(middleware.PermDeviceRead), devH.List)
			// 执行 dev.GET("/monitor", middleware.RequirePerm(middleware.PermDeviceRead), devH.Monitor)，完成对应的业务处理

			dev.GET("/monitor", middleware.RequirePerm(middleware.PermDeviceRead), devH.Monitor)
			// 执行 dev.POST("/:id/control", middleware.RequirePerm(middleware.PermDeviceCtl), devH.Control)，完成对应的业务处理

			dev.POST("/:id/control", middleware.RequirePerm(middleware.PermDeviceCtl), devH.Control)
		// 当前代码块结束

		}

		// 充值对账：读(网点管理员+运营)/执行(仅运营)
	// 定义并初始化局部变量 rec，取值为 admin.Group("/reconcile")

	rec := admin.Group("/reconcile")
		// 进入该代码块，开始执行其中的逻辑

		{
			// 执行 rec.GET("", middleware.RequirePerm(middleware.PermReconRead), recH.List)，完成对应的业务处理

			rec.GET("", middleware.RequirePerm(middleware.PermReconRead), recH.List)
			// 执行 rec.POST("/run", middleware.RequirePerm(middleware.PermReconRun), recH.Run)，完成对应的业务处理

			rec.POST("/run", middleware.RequirePerm(middleware.PermReconRun), recH.Run)
		// 当前代码块结束

		}

		// IoT 中台：总览/指令/OTA（读走 device:read，写走 ota:manage）@ticket B端-IoT中台
		// 定义并初始化局部变量 iotH，取值为 handler.NewIot(serviceIot(st))

		iotH := handler.NewIot(serviceIot(st))
		// 定义并初始化局部变量 iot，取值为 admin.Group("/iot")

		iot := admin.Group("/iot")
		// 进入该代码块，开始执行其中的逻辑

		{
			// 中台总览：设备状态分布/EMQX 健康/Topic 与 ACL 规范
			// 执行 iot.GET("/overview", middleware.RequirePerm(middleware.PermDeviceRead), iotH.Overview)，完成对应的业务处理

			iot.GET("/overview", middleware.RequirePerm(middleware.PermDeviceRead), iotH.Overview)
			// 下发控制指令（平台运营/网点管理员）
			// 执行 iot.POST("/devices/:id/command", middleware.RequirePerm(middleware.PermOtaManage), iotH.Dispatch)，完成对应的业务处理

			iot.POST("/devices/:id/command", middleware.RequirePerm(middleware.PermOtaManage), iotH.Dispatch)
			// 指令响应记录
			// 执行 iot.GET("/commands", middleware.RequirePerm(middleware.PermDeviceRead), iotH.CommandRecords)，完成对应的业务处理

			iot.GET("/commands", middleware.RequirePerm(middleware.PermDeviceRead), iotH.CommandRecords)
			// 新建 OTA 升级任务
			// 执行 iot.POST("/devices/:id/ota", middleware.RequirePerm(middleware.PermOtaManage), iotH.OtaCreate)，完成对应的业务处理

			iot.POST("/devices/:id/ota", middleware.RequirePerm(middleware.PermOtaManage), iotH.OtaCreate)
			// OTA 任务列表
			// 执行 iot.GET("/ota-tasks", middleware.RequirePerm(middleware.PermDeviceRead), iotH.OtaTasks)，完成对应的业务处理

			iot.GET("/ota-tasks", middleware.RequirePerm(middleware.PermDeviceRead), iotH.OtaTasks)
		// 当前代码块结束

		}

		// 告警处理：读/处理
		// 定义并初始化局部变量 alarm，取值为 admin.Group("/alarms")

		alarm := admin.Group("/alarms")
		// 进入该代码块，开始执行其中的逻辑

		{
			// 执行 alarm.GET("", middleware.RequirePerm(middleware.PermAlarmRead), alarmH.List)，完成对应的业务处理

			alarm.GET("", middleware.RequirePerm(middleware.PermAlarmRead), alarmH.List)
			// 执行 alarm.POST("/:id/handle", middleware.RequirePerm(middleware.PermAlarmHandle), alarmH.Handle)，完成对应的业务处理

			alarm.POST("/:id/handle", middleware.RequirePerm(middleware.PermAlarmHandle), alarmH.Handle)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 r 给调用方，结束当前函数

	return r
// 当前代码块结束

}
