// @ticket IoT-MQTT设备接入：设备 REST 接口（设备凭证/EMQX Hook/Webhook/指令/OTA/事件查询）
// 路由全部使用静态路径段，避免 gin 静态/通配符同级冲突
package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// DeviceAuth 设备请求 MQTT 凭证：POST /api/v1/devices/auth（设备侧调用，无需用户 JWT）@ticket IoT-MQTT设备接入
// 定义函数 DeviceAuth：封装一段可复用的处理逻辑（参数：c *gin.Context）

func DeviceAuth(c *gin.Context) {
	// 定义并初始化局部变量 req，取值为 map 结构并绑定 JSON 入参

	var req struct {
		// 结构体字段 DeviceSN，类型为 string，附带标签 `json:"device_sn"`；设备序列号

		DeviceSN        string `json:"device_sn"`
		// 结构体字段 DeviceSecret，类型为 string，附带标签 `json:"device_secret"`；出厂密钥

		DeviceSecret    string `json:"device_secret"`
		// 结构体字段 FirmwareVersion，类型为 string，附带标签 `json:"firmware_version"`；固件版本（可选）

		FirmwareVersion string `json:"firmware_version"`
	// 当前代码块结束

	}
	// 条件判断：err := c.ShouldBindJSON(&req); err != nil 时执行下面逻辑

	if err := c.ShouldBindJSON(&req); err != nil {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 调用业务层签发凭证（客户端 IP 用于审计与设备档案）@ticket IoT-MQTT设备接入
	// 定义并初始化局部变量 data, code，取值为 service.DeviceAuth(req.DeviceSN, req.DeviceSecret, req.FirmwareVersion, c.ClientIP())

	data, code := service.DeviceAuth(req.DeviceSN, req.DeviceSecret, req.FirmwareVersion, c.ClientIP())
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// EMQXAuthHook EMQX HTTP 认证 Hook：POST /api/v1/devices/emqx/auth（允许 200 / 拒绝 401）@ticket IoT-MQTT设备接入
// 兼容 JSON 与 form 两种提交方式（方案 2.2 步骤 8~11）
// 定义函数 EMQXAuthHook：封装一段可复用的处理逻辑（参数：c *gin.Context）

func EMQXAuthHook(c *gin.Context) {
	// 定义并初始化局部变量 username, password，取值为空串

	username, password := "", ""
	// 条件判断：c.ContentType() 含 application/json 时执行下面逻辑

	if c.ContentType() == "application/json" {
		// 定义并初始化局部变量 req，取值为 map 结构

		var req struct {
			// 结构体字段 Username，类型为 string，附带标签 `json:"username"`

			Username string `json:"username"`
			// 结构体字段 Password，类型为 string，附带标签 `json:"password"`

			Password string `json:"password"`
		// 当前代码块结束

		}
		// 条件判断：err := c.ShouldBindJSON(&req); err == nil 时执行下面逻辑

		if err := c.ShouldBindJSON(&req); err == nil {
			// 为 username, password 赋值：req.Username, req.Password

			username, password = req.Username, req.Password
		// 当前代码块结束

		}
	// 上述条件不成立时，执行 else 分支逻辑

	} else {
		// 为 username, password 赋值：c.PostForm("username"), c.PostForm("password")

		username, password = c.PostForm("username"), c.PostForm("password")
	// 当前代码块结束

	}
	// 校验 Token（方案 2.2：查询 Token 有效性并返回认证结果）
	// 条件判断：!service.VerifyEMQXToken(username, password) 时执行下面逻辑

	if !service.VerifyEMQXToken(username, password) {
		// 返回 401 拒绝接入，报文含 result=deny 便于排查
		// 执行 c.JSON(http.StatusUnauthorized, gin.H{"result": "deny", "message": "invalid device token"})，完成对应的业务处理

		c.JSON(http.StatusUnauthorized, gin.H{"result": "deny", "message": "invalid device token"})
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 返回 200 允许接入
	// 执行 c.JSON(http.StatusOK, gin.H{"result": "allow"})，完成对应的业务处理

	c.JSON(http.StatusOK, gin.H{"result": "allow"})
// 当前代码块结束

}

// EMQXWebhook EMQX 上下线 Webhook：POST /api/v1/devices/emqx/webhook（client_connected/client_disconnected）@ticket IoT-MQTT设备接入
// client_id 形如 device_{设备ID}（方案 2.2 连接参数）
// 定义函数 EMQXWebhook：封装一段可复用的处理逻辑（参数：c *gin.Context）

func EMQXWebhook(c *gin.Context) {
	// 定义并初始化局部变量 req，取值为 map 结构并绑定 JSON 入参

	var req struct {
		// 结构体字段 Event，类型为 string，附带标签 `json:"event"`

		Event    string `json:"event"`
		// 结构体字段 ClientID，类型为 string，附带标签 `json:"clientid"`

		ClientID string `json:"clientid"`
		// 结构体字段 IPAddress，类型为 string，附带标签 `json:"ip_address"`

		IPAddress string `json:"ip_address"`
	// 当前代码块结束

	}
	// 条件判断：err := c.ShouldBindJSON(&req); err != nil 时执行下面逻辑

	if err := c.ShouldBindJSON(&req); err != nil {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 解析 device_{id} 后缀得到设备 ID
	// 定义并初始化局部变量 deviceID, ok，取值为 parseDeviceClientID(req.ClientID)

	deviceID, ok := parseDeviceClientID(req.ClientID)
	// 条件判断：!ok 时执行下面逻辑

	if !ok {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 多分支判断：根据 req.Event 的值选择分支

	switch req.Event {
	// 当取值匹配 client_connected 时进入该分支（设备上线）
	case "client_connected":
		// 执行 service.HandleConnected(deviceID, req.IPAddress)，完成对应的业务处理

		service.HandleConnected(deviceID, req.IPAddress)
	// 当取值匹配 client_disconnected 时进入该分支（设备离线）
	case "client_disconnected":
		// 执行 service.HandleDisconnected(deviceID)，完成对应的业务处理

		service.HandleDisconnected(deviceID)
	// 当前代码块结束

	}
	// 执行 pkg.OK(c, gin.H{"ok": true})，完成对应的业务处理

	pkg.OK(c, gin.H{"ok": true})
// 当前代码块结束

}

// parseDeviceClientID 解析 device_{设备ID} 形式的 MQTT Client ID
// 定义函数 parseDeviceClientID：封装一段可复用的处理逻辑（参数：clientID string）(uint, bool)

func parseDeviceClientID(clientID string) (uint, bool) {
	// 前缀与长度校验
	// 条件判断：len(clientID) < 8 || clientID[:7] != "device_" 时执行下面逻辑

	if len(clientID) < 8 || clientID[:7] != "device_" {
		// 返回 0, false 给调用方，结束当前函数

		return 0, false
	// 当前代码块结束

	}
	// 解析数字后缀
	// 定义并初始化局部变量 n, err，取值为 strconv.ParseUint(clientID[7:], 10, 64)

	n, err := strconv.ParseUint(clientID[7:], 10, 64)
	// 条件判断：err != nil || n == 0 时执行下面逻辑

	if err != nil || n == 0 {
		// 返回 0, false 给调用方，结束当前函数

		return 0, false
	// 当前代码块结束

	}
	// 返回 uint(n), true 给调用方，结束当前函数

	return uint(n), true
// 当前代码块结束

}

// DeviceCommandH 下发设备指令：POST /api/v1/devices/command（JWT 鉴权；柜门动作接口须鉴权+审计）@ticket IoT-MQTT设备接入
// 定义函数 DeviceCommandH：封装一段可复用的处理逻辑（参数：c *gin.Context）

func DeviceCommandH(c *gin.Context) {
	// 定义并初始化局部变量 req，取值为 map 结构并绑定 JSON 入参

	var req struct {
		// 结构体字段 DeviceSN，类型为 string，附带标签 `json:"device_sn"`；目标设备序列号

		DeviceSN string `json:"device_sn"`
		// 结构体字段 CmdType，类型为 string，附带标签 `json:"cmd_type"`；open_door/close_door/reboot/config

		CmdType  string `json:"cmd_type"`
		// 结构体字段 LockerNo，类型为 string，附带标签 `json:"locker_no"`；柜格编号（开门类指令）

		LockerNo string `json:"locker_no"`
		// 结构体字段 OrderNo，类型为 string，附带标签 `json:"order_no"`；关联订单号

		OrderNo  string `json:"order_no"`
		// 结构体字段 Extra，类型为 map[string]interface{}，附带标签 `json:"extra"`；扩展字段（如 price_config）

		Extra    map[string]interface{} `json:"extra"`
	// 当前代码块结束

	}
	// 条件判断：err := c.ShouldBindJSON(&req); err != nil 时执行下面逻辑

	if err := c.ShouldBindJSON(&req); err != nil {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 指令类型白名单校验（方案 3.3 下行指令集合）@ticket IoT-MQTT设备接入
	// 定义并初始化局部变量 valid，取值为 map[string]bool{...}

	valid := map[string]bool{"open_door": true, "close_door": true, "reboot": true, "config": true}
	// 条件判断：!valid[req.CmdType] 时执行下面逻辑

	if !valid[req.CmdType] {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 触发人取当前登录用户（审计）@ticket IoT-MQTT设备接入
	// 定义并初始化局部变量 data, code，取值为 service.DispatchCommand(req.DeviceSN, req.CmdType, req.LockerNo, req.OrderNo, currentUserName(c), req.Extra)

	data, code := service.DispatchCommand(req.DeviceSN, req.CmdType, req.LockerNo, req.OrderNo, currentUserName(c), req.Extra)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// currentUserName 取当前登录用户标识（审计用；未登录返回 device-console）
// 定义函数 currentUserName：封装一段可复用的处理逻辑（参数：c *gin.Context）string

func currentUserName(c *gin.Context) string {
	// 已登录用户返回 user_{uid}（Claims 仅含 uid，无用户名字段）@ticket IoT-MQTT设备接入
	// 条件判断：uid := middleware.UID(c); uid > 0 时执行下面逻辑

	if uid := middleware.UID(c); uid > 0 {
		// 返回 fmt.Sprintf("user_%d", uid) 给调用方，结束当前函数

		return fmt.Sprintf("user_%d", uid)
	// 当前代码块结束

	}
	// 返回 "device-console" 给调用方，结束当前函数

	return "device-console"
// 当前代码块结束

}

// queryUint 解析 query 数字参数，缺省返回 0
// 定义函数 queryUint：封装一段可复用的处理逻辑（参数：c *gin.Context, key string）uint

func queryUint(c *gin.Context, key string) uint {
	// 定义并初始化局部变量 n, err，取值为 strconv.ParseUint(c.Query(key), 10, 64)

	n, err := strconv.ParseUint(c.Query(key), 10, 64)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 0 给调用方，结束当前函数

		return 0
	// 当前代码块结束

	}
	// 返回 uint(n) 给调用方，结束当前函数

	return uint(n)
// 当前代码块结束

}

// DeviceCommandsH 指令记录查询：GET /api/v1/devices/commands?device_id=&page=&page_size=（JWT 鉴权）@ticket IoT-MQTT设备接入
// 定义函数 DeviceCommandsH：封装一段可复用的处理逻辑（参数：c *gin.Context）

func DeviceCommandsH(c *gin.Context) {
	// 定义并初始化局部变量 deviceID，取值为 queryUint(c, "device_id")

	deviceID := queryUint(c, "device_id")
	// 定义并初始化局部变量 data, code，取值为 service.CommandRecords(deviceID, int(queryUint(c, "page")), int(queryUint(c, "page_size")))

	data, code := service.CommandRecords(deviceID, int(queryUint(c, "page")), int(queryUint(c, "page_size")))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// DeviceEventsH 设备上行事件查询：GET /api/v1/devices/events?device_id=&kind=&page=&page_size=（JWT 鉴权）@ticket IoT-MQTT设备接入
// 定义函数 DeviceEventsH：封装一段可复用的处理逻辑（参数：c *gin.Context）

func DeviceEventsH(c *gin.Context) {
	// 定义并初始化局部变量 deviceID，取值为 queryUint(c, "device_id")

	deviceID := queryUint(c, "device_id")
	// 定义并初始化局部变量 kind，取值为 c.Query("kind")

	kind := c.Query("kind")
	// 定义并初始化局部变量 data, code，取值为 service.EventRecords(deviceID, kind, int(queryUint(c, "page")), int(queryUint(c, "page_size")))

	data, code := service.EventRecords(deviceID, kind, int(queryUint(c, "page")), int(queryUint(c, "page_size")))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// DeviceStatusH 设备状态列表：GET /api/v1/devices/status?site_id=（JWT 鉴权；密钥/IP 脱敏由模型标签控制）@ticket IoT-MQTT设备接入
// 定义函数 DeviceStatusH：封装一段可复用的处理逻辑（参数：c *gin.Context）

func DeviceStatusH(c *gin.Context) {
	// 定义并初始化局部变量 data, code，取值为 service.DeviceStatusList(queryUint(c, "site_id"))

	data, code := service.DeviceStatusList(queryUint(c, "site_id"))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// DeviceOtaCreateH 创建 OTA 升级任务：POST /api/v1/devices/ota（JWT 鉴权 + 审计）@ticket IoT-MQTT设备接入
// 定义函数 DeviceOtaCreateH：封装一段可复用的处理逻辑（参数：c *gin.Context）

func DeviceOtaCreateH(c *gin.Context) {
	// 定义并初始化局部变量 req，取值为 map 结构并绑定 JSON 入参

	var req struct {
		// 结构体字段 DeviceSN，类型为 string，附带标签 `json:"device_sn"`

		DeviceSN string `json:"device_sn"`
		// 结构体字段 FirmwareVersion，类型为 string，附带标签 `json:"firmware_version"`

		FirmwareVersion string `json:"firmware_version"`
		// 结构体字段 DownloadURL，类型为 string，附带标签 `json:"download_url"`

		DownloadURL     string `json:"download_url"`
		// 结构体字段 MD5，类型为 string，附带标签 `json:"md5"`

		MD5             string `json:"md5"`
		// 结构体字段 FileSize，类型为 int64，附带标签 `json:"file_size"`

		FileSize        int64  `json:"file_size"`
		// 结构体字段 ForceUpgrade，类型为 bool，附带标签 `json:"force_upgrade"`

		ForceUpgrade    bool   `json:"force_upgrade"`
	// 当前代码块结束

	}
	// 条件判断：err := c.ShouldBindJSON(&req); err != nil 时执行下面逻辑

	if err := c.ShouldBindJSON(&req); err != nil {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 data, code，取值为 service.CreateOta(req.DeviceSN, req.FirmwareVersion, req.DownloadURL, req.MD5, req.FileSize, req.ForceUpgrade)

	data, code := service.CreateOta(req.DeviceSN, req.FirmwareVersion, req.DownloadURL, req.MD5, req.FileSize, req.ForceUpgrade)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// DeviceOtaListH OTA 任务列表：GET /api/v1/devices/ota/list?device_id=&page=&page_size=（JWT 鉴权）@ticket IoT-MQTT设备接入
// 定义函数 DeviceOtaListH：封装一段可复用的处理逻辑（参数：c *gin.Context）

func DeviceOtaListH(c *gin.Context) {
	// 定义并初始化局部变量 deviceID，取值为 queryUint(c, "device_id")

	deviceID := queryUint(c, "device_id")
	// 定义并初始化局部变量 data, code，取值为 service.OtaRecords(deviceID, int(queryUint(c, "page")), int(queryUint(c, "page_size")))

	data, code := service.OtaRecords(deviceID, int(queryUint(c, "page")), int(queryUint(c, "page_size")))
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}
