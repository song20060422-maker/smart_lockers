// @ticket 03-首页：定位/城市/寄存点/指南 HTTP handler
// @proto 原型对齐：v-sites 网点柜组监控(NearbySites)；见 AGENT.MD 2.1
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// Locate @ticket 03-首页：POST /api/v1/home/locate
// 定义函数 Locate：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Locate(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 Lat          float64 `json:"lat"`，完成对应的业务处理

		Lat          float64 `json:"lat"`
		// 执行 Lng          float64 `json:"lng"`，完成对应的业务处理

		Lng          float64 `json:"lng"`
		// 执行 IP           string  `json:"ip"`，完成对应的业务处理

		IP           string  `json:"ip"`
		// 执行 FallbackCity string  `json:"fallback_city"`，完成对应的业务处理

		FallbackCity string  `json:"fallback_city"`
	// 当前代码块结束

	}
	// 条件判断：err := c.ShouldBindJSON(&req); err != nil && c.Request.ContentLength > 0 时执行下面逻辑

	if err := c.ShouldBindJSON(&req); err != nil && c.Request.ContentLength > 0 {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 条件判断：req.FallbackCity == "" 时执行下面逻辑

	if req.FallbackCity == "" {
		// 为 req.FallbackCity 赋值："郑州"

		req.FallbackCity = "郑州"
	// 当前代码块结束

	}
	// 定义并初始化局部变量 data, code，取值为 service.Locate(req.Lat, req.Lng, req.IP, req.FallbackCity)

	data, code := service.Locate(req.Lat, req.Lng, req.IP, req.FallbackCity)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// SearchCities @ticket 03-首页：GET /api/v1/home/cities?keyword=
// 定义函数 SearchCities：封装一段可复用的处理逻辑（参数：c *gin.Context）

func SearchCities(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.SearchCities(c.Query("keyword"))

	data, code := service.SearchCities(c.Query("keyword"))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// NearbySites @ticket 03-首页：GET /api/v1/home/sites?lat=&lng=&page=&page_size=
// 定义函数 NearbySites：封装一段可复用的处理逻辑（参数：c *gin.Context）

func NearbySites(c *gin.Context) {
	// 定义并初始化局部变量 lat, _，取值为 strconv.ParseFloat(c.Query("lat"), 64)

	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	// 定义并初始化局部变量 lng, _，取值为 strconv.ParseFloat(c.Query("lng"), 64)

	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	// 定义并初始化局部变量 page, _，取值为 strconv.Atoi(c.DefaultQuery("page", "1"))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// 定义并初始化局部变量 pageSize, _，取值为 strconv.Atoi(c.DefaultQuery("page_size", "20"))

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	// 定义并初始化局部变量 data, code，取值为 service.NearbySites(lat, lng, page, pageSize)

	data, code := service.NearbySites(lat, lng, page, pageSize)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// GetGuide @ticket 03-首页：GET /api/v1/home/guide?version=&lang=
// 定义函数 GetGuide：封装一段可复用的处理逻辑（参数：c *gin.Context）

func GetGuide(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.GetGuide(c.Query("version"), c.Query("lang"))

	data, code := service.GetGuide(c.Query("version"), c.Query("lang"))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}
