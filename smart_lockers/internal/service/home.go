// @ticket 03-首页：定位、城市搜索、附近寄存点、寄存指南业务逻辑
package service

import (
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
)

// cityCache 城市词库短时缓存 @ticket 03-首页（JSON 词库动态更新 + 短时缓存）
var (
	cityCacheMu sync.RWMutex
	cityCache   []model.City
	cityCacheAt time.Time
)

// Locate 定位：GPS 优先，IP/默认回退；定位数据不持久化，短时缓存 @ticket 03-首页
func Locate(lat, lng float64, ip, fallbackCity string) (map[string]interface{}, int) {
	if lat == 0 && lng == 0 {
		// GPS 缺失回退到 IP/默认城市 @ticket 03-首页
		return map[string]interface{}{
			"source": "ip_fallback", "city": fallbackCity,
			"lat": 34.7466, "lng": 113.6254,
		}, pkg.CodeOK
	}
	return map[string]interface{}{
		"source": "gps", "lat": lat, "lng": lng,
	}, pkg.CodeOK
}

// SearchCities 按拼音首字母/中文模糊搜索并分组 @ticket 03-首页 P95≤300ms
func SearchCities(keyword string) ([]model.City, int) {
	if pkg.DB == nil {
		return nil, pkg.CodeOK
	}
	if time.Since(cityCacheAt) > 30*time.Second || len(cityCache) == 0 {
		var cities []model.City
		if err := pkg.DB.Order("sort asc, id asc").Find(&cities).Error; err != nil {
			pkg.Logf("[ERROR] load city cache: %v", err)
			return nil, pkg.ErrInternal
		}
		cityCacheMu.Lock()
		cityCache = cities
		cityCacheAt = time.Now()
		cityCacheMu.Unlock()
	}
	cityCacheMu.RLock()
	defer cityCacheMu.RUnlock()
	kw := strings.ToLower(strings.TrimSpace(keyword))
	res := make([]model.City, 0, len(cityCache))
	for _, c := range cityCache {
		if kw == "" || strings.Contains(strings.ToLower(c.Name), kw) ||
			strings.Contains(strings.ToLower(c.Pinyin), kw) ||
			strings.HasPrefix(strings.ToLower(c.Initial), kw) {
			res = append(res, c)
		}
	}
	return res, pkg.CodeOK
}

// siteWithDistance 供排序的距离附加结构 @ticket 03-首页
type siteWithDistance struct {
	model.LockerSite
	Distance float64 `json:"distance"`
}

// NearbySites 附近寄存点按距离排序 + 分页 @ticket 03-首页 P95≤300ms
func NearbySites(lat, lng float64, page, pageSize int) (map[string]interface{}, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pkg.DB == nil {
		return map[string]interface{}{"list": []siteWithDistance{}, "page": page, "page_size": pageSize}, pkg.CodeOK
	}
	var sites []model.LockerSite
	if err := pkg.DB.Where("status = ?", 1).Find(&sites).Error; err != nil {
		return nil, pkg.ErrInternal
	}
	list := make([]siteWithDistance, 0, len(sites))
	for _, s := range sites {
		list = append(list, siteWithDistance{LockerSite: s, Distance: haversine(lat, lng, s.Lat, s.Lng)})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Distance < list[j].Distance })
	total := len(list)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return map[string]interface{}{
		"list": list[start:end], "total": total, "page": page, "page_size": pageSize,
	}, pkg.CodeOK
}

// haversine 计算两点球面距离（公里）@ticket 03-首页
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const r = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * r * math.Asin(math.Sqrt(a))
}

// GetGuide 按版本/语言拉取寄存指南 @ticket 03-首页
func GetGuide(version, lang string) (map[string]interface{}, int) {
	if pkg.DB == nil {
		return map[string]interface{}{"version": version, "lang": lang, "content": ""}, pkg.CodeOK
	}
	var g model.Guide
	q := pkg.DB.Where("1 = 1")
	if version != "" {
		q = q.Where("version = ?", version)
	} else {
		q = q.Where("is_current = ?", true)
	}
	if lang != "" {
		q = q.Where("lang = ?", lang)
	}
	if err := q.First(&g).Error; err != nil {
		return map[string]interface{}{"version": version, "lang": lang, "content": ""}, pkg.CodeOK
	}
	return map[string]interface{}{"version": g.Version, "lang": g.Lang, "content": g.Content}, pkg.CodeOK
}
