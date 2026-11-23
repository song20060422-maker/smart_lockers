// @ticket 07-数字人Coze：本地演示模式 WebSocket 后端 handler；握手读 token，收 user 消息回单条 assistant 文本帧
// 定义包 handler：集中存放 HTTP 与 WebSocket 路由的处理函数

package handler

import (
	// 引入 encoding/json：解析上行消息与下行回复序列化

	"encoding/json"
	// 引入 net/http：Upgrader 的握手请求类型与状态码

	"net/http"

	"github.com/gin-gonic/gin"
	// 引入 gorilla/websocket：实现 RFC6455 服务端升级与文本帧收发

	"github.com/gorilla/websocket"

	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// wsUpgrader 配置 WebSocket 升级器，本地演示模式放行所有 Origin
// 定义变量 wsUpgrader：用于把 HTTP 请求升级为 WebSocket 连接

var wsUpgrader = websocket.Upgrader{
	// 设置读缓冲区大小为 1024 字节

	ReadBufferSize: 1024,
	// 设置写缓冲区大小为 1024 字节

	WriteBufferSize: 1024,
	// 允许所有来源的握手请求，本地演示模式下前后端分离开发需要

	CheckOrigin: func(r *http.Request) bool { return true },
	// 当前代码块结束

}

// userWSMessage 严格对齐数字人前端 WebSocketClient.send 发出的上行消息结构
// 定义结构体 userWSMessage：用于解析前端发来的 user 消息

type userWSMessage struct {
	// 结构体字段 Role，取值为 user，附带标签 `json:"role"`

	Role string `json:"role"`
	// 结构体字段 Content，用户文本内容，附带标签 `json:"content"`

	Content string `json:"content"`
	// 结构体字段 Type，取值为 text 或 voice，附带标签 `json:"type"`

	Type string `json:"type"`
	// 结构体字段 Timestamp，毫秒时间戳，附带标签 `json:"timestamp"`

	Timestamp int64 `json:"timestamp"`
	// 当前代码块结束

}

// DigitalHumanWS 处理 /ws/digital-human 路由：GET 升级为 WebSocket，循环收发文本帧
// 定义函数 DigitalHumanWS：封装一段可复用的处理逻辑（参数：c *gin.Context）

func DigitalHumanWS(c *gin.Context) {
	// 从握手 query 读取 token，本地演示模式仅记录不强制校验

	token := c.Query("token")
	// 从握手 query 读取角色，决定本连接使用哪一套人设与专属问答 @ticket 07

	role := c.Query("role")
	// 条件判断：角色不在三业务角色白名单内时回退到智能客服专员

	if role != "sender" && role != "receiver" && role != "service" {
		role = "service"
	}
	// 记录握手日志，便于排查前端是否携带 token 与角色

	service.LogDemo("[digital-human] ws handshake token_present=%v role=%s", token != "", role)
	// 条件判断：token 为空时仍放行，但打印提示便于演示排查

	if token == "" {
		// 记录无 token 的提示日志

		service.LogDemo("[digital-human] handshake without token, demo mode allows anonymous")
		// 当前代码块结束

	}
	// 将 HTTP 连接升级为 WebSocket 连接

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	// 条件判断：升级失败时直接返回，框架已响应错误

	if err != nil {
		// 记录升级失败日志

		service.LogDemo("[digital-human] upgrade failed: %v", err)
		// 返回  给调用方，结束当前函数

		return
		// 当前代码块结束

	}
	// 函数返回时关闭连接，释放资源

	defer conn.Close()
	// 记录连接建立日志
	service.LogDemo("[digital-human] client connected role=%s", role)
	// 连接建立即主动推送当前角色的欢迎语，走与问答一致的口型/动画链路 @ticket 07
	welcomePayload, werr := json.Marshal(service.BuildWelcome(role))
	// 条件判断：欢迎语序列化或写入失败仅记录日志，不阻断后续问答循环
	if werr == nil {
		if err := conn.WriteMessage(websocket.TextMessage, welcomePayload); err != nil {
			service.LogDemo("[digital-human] write welcome failed: %v", err)
		}
	} else {
		service.LogDemo("[digital-human] marshal welcome failed: %v", werr)
	}
	// 进入无限循环，持续读取客户端消息
	for {
		// 读取一条 WebSocket 消息，msgType 为帧类型，data 为载荷

		msgType, data, err := conn.ReadMessage()
		// 条件判断：读取失败时跳出循环，由 defer 关闭连接

		if err != nil {
			// 记录断开原因日志

			service.LogDemo("[digital-human] client disconnected: %v", err)
			// 跳出循环

			break
			// 当前代码块结束

		}
		// 条件判断：仅处理文本帧，其它帧类型忽略

		if msgType != websocket.TextMessage {
			// 进入下一轮循环

			continue
			// 当前代码块结束

		}
		// 定义并初始化局部变量 msg，用于解析上行消息

		var msg userWSMessage
		// 条件判断：解析上行 JSON 失败时忽略该帧，保证无控制台错误

		if err := json.Unmarshal(data, &msg); err != nil {
			// 记录非法帧日志

			service.LogDemo("[digital-human] ignore invalid frame: %v", err)
			// 进入下一轮循环

			continue
			// 当前代码块结束

		}
		// 条件判断：仅处理 role 为 user 且内容非空的消息

		if msg.Role != "user" || msg.Content == "" {
			// 进入下一轮循环

			continue
			// 当前代码块结束

		}
		// 调用本地问答引擎构造严格对齐前端协议的回复，传入本连接角色实现内容差异化 @ticket 07
		reply := service.BuildReply(msg.Content, role)
		// 将回复序列化为单条 JSON 文本帧

		payload, err := json.Marshal(reply)
		// 条件判断：序列化失败时记录日志并跳过

		if err != nil {
			// 记录序列化失败日志

			service.LogDemo("[digital-human] marshal reply failed: %v", err)
			// 进入下一轮循环

			continue
			// 当前代码块结束

		}
		// 条件判断：写入文本帧失败时记录日志并跳出循环

		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			// 记录写入失败日志

			service.LogDemo("[digital-human] write failed: %v", err)
			// 跳出循环

			break
			// 当前代码块结束

		}
		// 当前代码块结束

	}
	// 当前代码块结束

}

// digitalHumanLoginReq 数字人账号密码登录请求体，字段契约与前端 AccountLoginPage 对齐 @ticket 07-数字人Coze
// 定义结构体 digitalHumanLoginReq：承载 username/password 两个字段
type digitalHumanLoginReq struct {
	// 结构体字段 Username，登录账号
	Username string `json:"username"`
	// 结构体字段 Password，明文密码（传输层依赖 HTTPS，服务端 bcrypt 比对）
	Password string `json:"password"`
}

// DigitalHumanLogin 数字人账号密码登录 HTTP 接口：真实查库校验 @ticket 07-数字人Coze
// 定义函数 DigitalHumanLogin：解析请求 → 调用 service 查库校验 → respond 统一出参
func DigitalHumanLogin(c *gin.Context) {
	// 定义并初始化局部变量 req，承载请求体
	var req digitalHumanLoginReq
	// 绑定 JSON 请求体，参数错误直接返回业务码
	if err := c.ShouldBindJSON(&req); err != nil {
		// 返回参数错误响应
		respond(c, nil, pkg.ErrBad)
		// 结束当前函数
		return
	}
	// 调用 service 层查库校验账号密码
	data, code := service.DigitalHumanLogin(req.Username, req.Password)
	// 统一响应：成功返回 token/username/role，失败返回业务码
	respond(c, data, code)
}
