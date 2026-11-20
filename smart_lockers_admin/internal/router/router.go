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
func Setup(st *store.Store) *gin.Engine {
	gin.SetMode(config.Get().Server.Mode)
	r := gin.New()
	r.Use(middleware.Recovery(), middleware.RequestID(), middleware.Logger(), middleware.CORS())

	r.GET("/health", func(c *gin.Context) { pkg.OK(c, map[string]interface{}{"status": "up", "service": "smart_lockers_admin"}) })

	authH := &handler.AuthHandler{}
	siteH := handler.NewSite(serviceSite(st))
	devH := handler.NewDevice(serviceDevice(st))
	recH := handler.NewReconcile(serviceRecon(st))
	alarmH := handler.NewAlarm(serviceAlarm(st))

	v1 := r.Group("/api/v1")

	// 登录：签令牌
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authH.Login)
		auth.GET("/me", middleware.Auth(), authH.Me)
	}

	// 其余接口统一鉴权
	admin := v1.Group("/admin")
	admin.Use(middleware.Auth())
	{
		// 网点管理：读(3类角色)/写(网点管理员+运营)
		sites := admin.Group("/sites")
		{
			sites.GET("", middleware.RequirePerm(middleware.PermSiteRead), siteH.List)
			sites.POST("", middleware.RequirePerm(middleware.PermSiteWrite), siteH.Create)
			sites.POST("/:id/status", middleware.RequirePerm(middleware.PermSiteWrite), siteH.UpdateStatus)
		}

		// 设备监控：读/控制
		dev := admin.Group("/devices")
		{
			dev.GET("", middleware.RequirePerm(middleware.PermDeviceRead), devH.List)
			dev.GET("/monitor", middleware.RequirePerm(middleware.PermDeviceRead), devH.Monitor)
			dev.POST("/:id/control", middleware.RequirePerm(middleware.PermDeviceCtl), devH.Control)
		}

		// 充值对账：读(网点管理员+运营)/执行(仅运营)
		rec := admin.Group("/reconcile")
		{
			rec.GET("", middleware.RequirePerm(middleware.PermReconRead), recH.List)
			rec.POST("/run", middleware.RequirePerm(middleware.PermReconRun), recH.Run)
		}

		// 告警处理：读/处理
		alarm := admin.Group("/alarms")
		{
			alarm.GET("", middleware.RequirePerm(middleware.PermAlarmRead), alarmH.List)
			alarm.POST("/:id/handle", middleware.RequirePerm(middleware.PermAlarmHandle), alarmH.Handle)
		}
	}
	return r
}
