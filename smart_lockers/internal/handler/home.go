// @ticket 03-首页：定位/城市/寄存点/指南 HTTP handler
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// Locate @ticket 03-首页：POST /api/v1/home/locate
func Locate(c *gin.Context) {
	var req struct {
		Lat          float64 `json:"lat"`
		Lng          float64 `json:"lng"`
		IP           string  `json:"ip"`
		FallbackCity string  `json:"fallback_city"`
	}
	if err := c.ShouldBindJSON(&req); err != nil && c.Request.ContentLength > 0 {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	if req.FallbackCity == "" {
		req.FallbackCity = "郑州"
	}
	data, code := service.Locate(req.Lat, req.Lng, req.IP, req.FallbackCity)
	respond(c, data, code)
}

// SearchCities @ticket 03-首页：GET /api/v1/home/cities?keyword=
func SearchCities(c *gin.Context) {
	data, code := service.SearchCities(c.Query("keyword"))
	respond(c, data, code)
}

// NearbySites @ticket 03-首页：GET /api/v1/home/sites?lat=&lng=&page=&page_size=
func NearbySites(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	data, code := service.NearbySites(lat, lng, page, pageSize)
	respond(c, data, code)
}

// GetGuide @ticket 03-首页：GET /api/v1/home/guide?version=&lang=
func GetGuide(c *gin.Context) {
	data, code := service.GetGuide(c.Query("version"), c.Query("lang"))
	respond(c, data, code)
}
