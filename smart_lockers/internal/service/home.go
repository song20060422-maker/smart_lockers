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
// 开始定义一组包级/局部变量的 var 声明块

var (
	// 执行 cityCacheMu sync.RWMutex，完成对应的业务处理

	cityCacheMu sync.RWMutex
	// 执行 cityCache   []model.City，完成对应的业务处理

	cityCache   []model.City
	// 执行 cityCacheAt time.Time，完成对应的业务处理

	cityCacheAt time.Time
// 结束参数/表达式列表

)

// Locate 定位：GPS 优先，IP/默认回退；定位数据不持久化，短时缓存 @ticket 03-首页
// 定义函数 Locate：封装一段可复用的处理逻辑（参数：lat, lng float64, ip, fallbackCity string）

func Locate(lat, lng float64, ip, fallbackCity string) (map[string]interface{}, int) {
	// 条件判断：lat == 0 && lng == 0 时执行下面逻辑

	if lat == 0 && lng == 0 {
		// GPS 缺失回退到 IP/默认城市 @ticket 03-首页
		// 返回 map[string]interface{}{ 给调用方，结束当前函数

		return map[string]interface{}{
			// 执行 "source": "ip_fallback", "city": fallbackCity,，完成对应的业务处理

			"source": "ip_fallback", "city": fallbackCity,
			// 执行 "lat": 34.7466, "lng": 113.6254,，完成对应的业务处理

			"lat": 34.7466, "lng": 113.6254,
		// 执行 }, pkg.CodeOK，完成对应的业务处理

		}, pkg.CodeOK
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{ 给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "source": "gps", "lat": lat, "lng": lng,，完成对应的业务处理

		"source": "gps", "lat": lat, "lng": lng,
	// 执行 }, pkg.CodeOK，完成对应的业务处理

	}, pkg.CodeOK
// 当前代码块结束

}

// SearchCities 按拼音首字母/中文模糊搜索并分组 @ticket 03-首页 P95≤300ms
// 定义函数 SearchCities：封装一段可复用的处理逻辑（参数：keyword string）

func SearchCities(keyword string) ([]model.City, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 nil, pkg.CodeOK 给调用方，结束当前函数

		return nil, pkg.CodeOK
	// 当前代码块结束

	}
	// 条件判断：time.Since(cityCacheAt) > 30*time.Second || len(cityCache) == 0 时执行下面逻辑

	if time.Since(cityCacheAt) > 30*time.Second || len(cityCache) == 0 {
		// 声明变量 cities []model.City，用于保存运行过程中的数据

		var cities []model.City
		// 条件判断：err := pkg.DB.Order("sort asc, id asc").Find(&cities).Error; err != nil 时执行下面逻辑

		if err := pkg.DB.Order("sort asc, id asc").Find(&cities).Error; err != nil {
			// 执行 pkg.Logf("[ERROR] load city cache: %v", err)，完成对应的业务处理

			pkg.Logf("[ERROR] load city cache: %v", err)
			// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

			return nil, pkg.ErrInternal
		// 当前代码块结束

		}
		// 执行 cityCacheMu.Lock()，完成对应的业务处理

		cityCacheMu.Lock()
		// 为 cityCache 赋值：cities

		cityCache = cities
		// 为 cityCacheAt 赋值：time.Now()

		cityCacheAt = time.Now()
		// 执行 cityCacheMu.Unlock()，完成对应的业务处理

		cityCacheMu.Unlock()
	// 当前代码块结束

	}
	// 执行 cityCacheMu.RLock()，完成对应的业务处理

	cityCacheMu.RLock()
	// 延迟注册 cityCacheMu.RUnlock()，在当前函数返回前自动执行

	defer cityCacheMu.RUnlock()
	// 定义并初始化局部变量 kw，取值为 strings.ToLower(strings.TrimSpace(keyword))

	kw := strings.ToLower(strings.TrimSpace(keyword))
	// 定义并初始化局部变量 res，取值为 make([]model.City, 0, len(cityCache))

	res := make([]model.City, 0, len(cityCache))
	// 循环处理：for _, c := range cityCache，逐轮迭代

	for _, c := range cityCache {
		// 条件判断：kw == "" || strings.Contains(strings.ToLower(c.Name), kw) || 时执行下面逻辑

		if kw == "" || strings.Contains(strings.ToLower(c.Name), kw) ||
			// 执行 strings.Contains(strings.ToLower(c.Pinyin), kw) ||，完成对应的业务处理

			strings.Contains(strings.ToLower(c.Pinyin), kw) ||
			// 执行 strings.HasPrefix(strings.ToLower(c.Initial), kw) {，完成对应的业务处理

			strings.HasPrefix(strings.ToLower(c.Initial), kw) {
			// 为 res 赋值：append(res, c)

			res = append(res, c)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 res, pkg.CodeOK 给调用方，结束当前函数

	return res, pkg.CodeOK
// 当前代码块结束

}

// siteWithDistance 供排序的距离附加结构 @ticket 03-首页
// 定义结构体 siteWithDistance：用于承载一组相关的字段数据

type siteWithDistance struct {
	// 执行 model.LockerSite，完成对应的业务处理

	model.LockerSite
	// 结构体字段 Distance，类型为 float64，附带标签 `json:"distance"`

	Distance float64 `json:"distance"`
// 当前代码块结束

}

// NearbySites 附近寄存点按距离排序 + 分页 @ticket 03-首页 P95≤300ms
// 定义函数 NearbySites：封装一段可复用的处理逻辑（参数：lat, lng float64, page, pageSize int）

func NearbySites(lat, lng float64, page, pageSize int) (map[string]interface{}, int) {
	// 条件判断：page <= 0 时执行下面逻辑

	if page <= 0 {
		// 为 page 赋值：1

		page = 1
	// 当前代码块结束

	}
	// 条件判断：pageSize <= 0 时执行下面逻辑

	if pageSize <= 0 {
		// 为 pageSize 赋值：20

		pageSize = 20
	// 当前代码块结束

	}
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 map[string]interface{}{"list": []siteWithDistance{}, "page": page, "page_size": pageSize}, pkg.CodeOK 给调用方，结束当前函数

		return map[string]interface{}{"list": []siteWithDistance{}, "page": page, "page_size": pageSize}, pkg.CodeOK
	// 当前代码块结束

	}
	// 声明变量 sites []model.LockerSite，用于保存运行过程中的数据

	var sites []model.LockerSite
	// 条件判断：err := pkg.DB.Where("status = ?", 1).Find(&sites).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("status = ?", 1).Find(&sites).Error; err != nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 定义并初始化局部变量 list，取值为 make([]siteWithDistance, 0, len(sites))

	list := make([]siteWithDistance, 0, len(sites))
	// 循环处理：for _, s := range sites，逐轮迭代

	for _, s := range sites {
		// 为 list 赋值：append(list, siteWithDistance{LockerSite: s, Distance: haversine(lat, lng, s.Lat, s.Lng)})

		list = append(list, siteWithDistance{LockerSite: s, Distance: haversine(lat, lng, s.Lat, s.Lng)})
	// 当前代码块结束

	}
	// 执行 sort.Slice(list, func(i, j int) bool { return list[i].Distance < list[j].Distance })，完成对应的业务处理

	sort.Slice(list, func(i, j int) bool { return list[i].Distance < list[j].Distance })
	// 定义并初始化局部变量 total，取值为 len(list)

	total := len(list)
	// 定义并初始化局部变量 start，取值为 (page - 1) * pageSize

	start := (page - 1) * pageSize
	// 条件判断：start > total 时执行下面逻辑

	if start > total {
		// 为 start 赋值：total

		start = total
	// 当前代码块结束

	}
	// 定义并初始化局部变量 end，取值为 start + pageSize

	end := start + pageSize
	// 条件判断：end > total 时执行下面逻辑

	if end > total {
		// 为 end 赋值：total

		end = total
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{ 给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "list": list[start:end], "total": total, "page": page, "page_size": pageSize,，完成对应的业务处理

		"list": list[start:end], "total": total, "page": page, "page_size": pageSize,
	// 执行 }, pkg.CodeOK，完成对应的业务处理

	}, pkg.CodeOK
// 当前代码块结束

}

// haversine 计算两点球面距离（公里）@ticket 03-首页
// 定义函数 haversine：封装一段可复用的处理逻辑（参数：lat1, lng1, lat2, lng2 float64）

func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	// 声明常量 r = 6371.0：固定不变的取值

	const r = 6371.0
	// 定义并初始化局部变量 dLat，取值为 (lat2 - lat1) * math.Pi / 180

	dLat := (lat2 - lat1) * math.Pi / 180
	// 定义并初始化局部变量 dLng，取值为 (lng2 - lng1) * math.Pi / 180

	dLng := (lng2 - lng1) * math.Pi / 180
	// 定义并初始化局部变量 a，取值为 math.Sin(dLat/2)*math.Sin(dLat/2) +

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		// 执行 math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLng/2)*math.Sin(dLng/2)，完成对应的业务处理

		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*math.Sin(dLng/2)*math.Sin(dLng/2)
	// 返回 2 * r * math.Asin(math.Sqrt(a)) 给调用方，结束当前函数

	return 2 * r * math.Asin(math.Sqrt(a))
// 当前代码块结束

}

// GetGuide 按版本/语言拉取寄存指南 @ticket 03-首页
// 定义函数 GetGuide：封装一段可复用的处理逻辑（参数：version, lang string）

func GetGuide(version, lang string) (map[string]interface{}, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 map[string]interface{}{"version": version, "lang": lang, "content": ""}, pkg.CodeOK 给调用方，结束当前函数

		return map[string]interface{}{"version": version, "lang": lang, "content": ""}, pkg.CodeOK
	// 当前代码块结束

	}
	// 声明变量 g model.Guide，用于保存运行过程中的数据

	var g model.Guide
	// 定义并初始化局部变量 q，取值为 pkg.DB.Where("1 = 1")

	q := pkg.DB.Where("1 = 1")
	// 条件判断：version != "" 时执行下面逻辑

	if version != "" {
		// 为 q 赋值：q.Where("version = ?", version)

		q = q.Where("version = ?", version)
	// 上述条件不成立时，执行 else 分支逻辑

	} else {
		// 为 q 赋值：q.Where("is_current = ?", true)

		q = q.Where("is_current = ?", true)
	// 当前代码块结束

	}
	// 条件判断：lang != "" 时执行下面逻辑

	if lang != "" {
		// 为 q 赋值：q.Where("lang = ?", lang)

		q = q.Where("lang = ?", lang)
	// 当前代码块结束

	}
	// 条件判断：err := q.First(&g).Error; err != nil 时执行下面逻辑

	if err := q.First(&g).Error; err != nil {
		// 返回 map[string]interface{}{"version": version, "lang": lang, "content": ""}, pkg.CodeOK 给调用方，结束当前函数

		return map[string]interface{}{"version": version, "lang": lang, "content": ""}, pkg.CodeOK
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"version": g.Version, "lang": g.Lang, "content": g.Content}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"version": g.Version, "lang": g.Lang, "content": g.Content}, pkg.CodeOK
// 当前代码块结束

}
