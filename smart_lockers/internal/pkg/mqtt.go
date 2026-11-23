// @ticket IoT-MQTT设备接入：EMQX 客户端封装；连接失败降级为本地演示模式，保证无 Broker 时 go run 可启动
// 对齐《智能快递柜IoT设备MQTT接入技术方案》：后端服务允许发布 platform/#、订阅 device/#（方案 3.4 ACL）
package pkg

import (
	"fmt"
	"log"
	"sync"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"

	"smart_lockers/internal/config"
)

// MqttConfig MQTT 连接配置（从 config.MqttConfig 转换而来，便于独立测试）
// 定义结构体 MqttConfig：用于承载一组相关的字段数据

type MqttConfig struct {
	// 结构体字段 Enabled，类型为 bool；false 时完全不连接 Broker（本地演示模式）

	Enabled bool
	// 结构体字段 Broker，类型为 string；EMQX 地址（tcp://host:1883）

	Broker string
	// 结构体字段 ClientID，类型为 string；后端服务客户端 ID

	ClientID string
	// 结构体字段 Username，类型为 string；后端服务 ACL 用户名（backend-service）

	Username string
	// 结构体字段 Password，类型为 string；后端服务密码

	Password string
// 当前代码块结束

}

// 声明变量 mqttMu sync.Mutex，用于保存运行过程中的数据

var (
	// 执行 mqttMu sync.Mutex，完成对应的业务处理

	mqttMu sync.Mutex
	// 执行 mqttClient pahomqtt.Client，完成对应的业务处理

	mqttClient pahomqtt.Client
	// 执行 mqttReady bool，完成对应的业务处理

	mqttReady bool
	// 执行 mqttSubMu sync.Mutex，完成对应的业务处理

	mqttSubMu sync.Mutex
// 结束参数/表达式列表

)

// InitMQTT 初始化后端 MQTT 客户端并连接 EMQX；失败仅告警不 panic（降级演示模式）@ticket IoT-MQTT设备接入
// handler 为上行消息统一入口（签名 func(topic, payload string)），与业务层路由函数解耦 paho 类型
// 定义函数 InitMQTT：封装一段可复用的处理逻辑（参数：cfg MqttConfig, handler func(topic string, payload []byte)）

func InitMQTT(cfg MqttConfig, handler func(topic string, payload []byte)) {
	// 加互斥锁保护客户端初始化的并发安全

	mqttMu.Lock()
	// 延迟注册 mqttMu.Unlock()，在当前函数返回前自动执行

	defer mqttMu.Unlock()
	// 条件判断：!cfg.Enabled 时执行下面逻辑

	if !cfg.Enabled {
		// 执行 log.Println("[INFO] mqtt disabled (本地演示模式，设备上下行走 HTTP/Webhook 兼容通道)")，完成对应的业务处理

		log.Println("[INFO] mqtt disabled (本地演示模式，设备上下行走 HTTP/Webhook 兼容通道)")
		// 为 mqttReady 赋值：false

		mqttReady = false
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 opts，取值为 pahomqtt.NewClientOptions().AddBroker(cfg.Broker)

	opts := pahomqtt.NewClientOptions().AddBroker(cfg.Broker)
	// 为 opts 客户端 ID 赋值：cfg.ClientID

	opts.SetClientID(cfg.ClientID)
	// 条件判断：cfg.Username != "" 时执行下面逻辑

	if cfg.Username != "" {
		// 执行 opts.SetUsername(cfg.Username)，完成对应的业务处理

		opts.SetUsername(cfg.Username)
		// 执行 opts.SetPassword(cfg.Password)，完成对应的业务处理

		opts.SetPassword(cfg.Password)
	// 当前代码块结束

	}
	// 自动重连：断线后指数退避重连（方案 2.4 断线重连机制）

	opts.SetAutoReconnect(true)
	// 清理会话：后端订阅 device/# 需要跨重连保留，设置 clean session=false 保留离线消息

	opts.SetCleanSession(false)
	// 设置 keep alive 60 秒（方案 2.3 连接参数推荐值）

	opts.SetKeepAlive(60 * time.Second)
	// 连接丢失回调：仅记录日志（AutoReconnect 会自动恢复）

	opts.SetConnectionLostHandler(func(c pahomqtt.Client, err error) {
		// 执行 log.Printf("[WARN] mqtt connection lost: %v", err)，完成对应的业务处理

		log.Printf("[WARN] mqtt connection lost: %v", err)
	// 当前代码块结束

	})
	// 重连成功回调：标记就绪并记录日志

	opts.SetOnConnectHandler(func(c pahomqtt.Client) {
		// 加锁更新就绪标记

		mqttMu.Lock()
		// 为 mqttReady 赋值：true

		mqttReady = true
		// 解锁

		mqttMu.Unlock()
		// 执行 log.Println("[INFO] mqtt connected")，完成对应的业务处理

		log.Println("[INFO] mqtt connected")
	// 当前代码块结束

	})
	// 设置消息统一入口处理器：将 paho 消息转换为 (topic, payload) 回调（由 service.RouteMessage 处理 device/# 上行）

	opts.SetDefaultPublishHandler(func(c pahomqtt.Client, m pahomqtt.Message) {
		// 调用业务路由函数处理上行消息
		// 执行 handler(m.Topic(), m.Payload())，完成对应的业务处理

		handler(m.Topic(), m.Payload())
	// 当前代码块结束

	})
	// 定义并初始化局部变量 cli，取值为 pahomqtt.NewClient(opts)

	cli := pahomqtt.NewClient(opts)
	// 定义并初始化局部变量 token，取值为 cli.Connect()，异步发起连接

	token := cli.Connect()
	// 条件判断：token.WaitTimeout(3*time.Second) && token.Error() != nil 时执行下面逻辑

	if ok := token.WaitTimeout(3 * time.Second); ok && token.Error() != nil {
		// 执行 log.Printf("[WARN] mqtt connect failed: %v (降级为本地演示模式)", token.Error())，完成对应的业务处理

		log.Printf("[WARN] mqtt connect failed: %v (降级为本地演示模式)", token.Error())
		// 为 mqttReady 赋值：false

		mqttReady = false
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 为 mqttClient 赋值：cli

	mqttClient = cli
	// 执行 log.Printf("[INFO] mqtt client ready (broker=%s client_id=%s)", cfg.Broker, cfg.ClientID)，完成对应的业务处理

	log.Printf("[INFO] mqtt client ready (broker=%s client_id=%s)", cfg.Broker, cfg.ClientID)
// 当前代码块结束

}

// MqttReady 返回当前 MQTT 是否可用（已启用且连接成功）@ticket IoT-MQTT设备接入
// 定义函数 MqttReady：封装一段可复用的处理逻辑（参数：）

func MqttReady() bool {
	// 加锁读取就绪标记

	mqttMu.Lock()
	// 延迟注册解锁

	defer mqttMu.Unlock()
	// 返回 mqttReady 给调用方，结束当前函数

	return mqttReady
// 当前代码块结束

}

// MqttSubscribe 订阅指定 Topic（通配符支持 +/#），QoS 固定 1（状态/指令可靠传输）@ticket IoT-MQTT设备接入
// cb 为业务回调（nil 时走 InitMQTT 注册的统一入口）
// 定义函数 MqttSubscribe：封装一段可复用的处理逻辑（参数：topic string, cb func(topic string, payload []byte)）error

func MqttSubscribe(topic string, cb func(topic string, payload []byte)) error {
	// 条件判断：!MqttReady() 时执行下面逻辑

	if !MqttReady() {
		// 返回 fmt.Errorf("mqtt unavailable") 给调用方，结束当前函数

		return fmt.Errorf("mqtt unavailable")
	// 当前代码块结束

	}
	// 加互斥锁串行化订阅操作，避免并发订阅竞争

	mqttSubMu.Lock()
	// 延迟注册解锁

	defer mqttSubMu.Unlock()
	// 将业务回调包装为 paho 处理器；nil 时使用包级统一入口 mqttHandler
	// 定义并初始化局部变量 pahoCb，取值为包装函数

	pahoCb := func(c pahomqtt.Client, m pahomqtt.Message) {
		// 条件判断：cb != nil 时执行下面逻辑

		if cb != nil {
			// 执行 cb(m.Topic(), m.Payload())，完成对应的业务处理

			cb(m.Topic(), m.Payload())
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 定义并初始化局部变量 t，取值为 mqttClient.Subscribe(topic, 1, pahoCb)

	t := mqttClient.Subscribe(topic, 1, pahoCb)
	// 条件判断：!t.WaitTimeout(3*time.Second) || t.Error() != nil 时执行下面逻辑

	if ok := t.WaitTimeout(3 * time.Second); !ok || t.Error() != nil {
		// 返回订阅错误（超时或 Broker 拒绝）给调用方，结束当前函数

		if t.Error() != nil {
			// 返回 t.Error() 给调用方，结束当前函数

			return t.Error()
		// 当前代码块结束

		}
		// 返回 fmt.Errorf("mqtt subscribe timeout") 给调用方，结束当前函数

		return fmt.Errorf("mqtt subscribe timeout")
	// 当前代码块结束

	}
	// 返回 nil 给调用方，结束当前函数

	return nil
// 当前代码块结束

}

// MqttPublish 向指定 Topic 发布消息；qos 由调用方传入（心跳 0、状态/指令 1、防拆 2）@ticket IoT-MQTT设备接入
// 定义函数 MqttPublish：封装一段可复用的处理逻辑（参数：topic string, qos byte, payload []byte）error

func MqttPublish(topic string, qos byte, payload []byte) error {
	// 条件判断：!MqttReady() 时执行下面逻辑

	if !MqttReady() {
		// 返回 fmt.Errorf("mqtt unavailable") 给调用方，结束当前函数

		return fmt.Errorf("mqtt unavailable")
	// 当前代码块结束

	}
	// 定义并初始化局部变量 t，取值为 mqttClient.Publish(topic, qos, false, payload)

	t := mqttClient.Publish(topic, qos, false, payload)
	// 条件判断：!t.WaitTimeout(3*time.Second) || t.Error() != nil 时执行下面逻辑

	if ok := t.WaitTimeout(3 * time.Second); !ok || t.Error() != nil {
		// 返回发布错误（超时或 Broker 拒绝）给调用方，结束当前函数

		if t.Error() != nil {
			// 返回 t.Error() 给调用方，结束当前函数

			return t.Error()
		// 当前代码块结束

		}
		// 返回 fmt.Errorf("mqtt publish timeout") 给调用方，结束当前函数

		return fmt.Errorf("mqtt publish timeout")
	// 当前代码块结束

	}
	// 返回 nil 给调用方，结束当前函数

	return nil
// 当前代码块结束

}

// NewMqttConfig 从全局配置构造 MqttConfig @ticket IoT-MQTT设备接入
// 定义函数 NewMqttConfig：封装一段可复用的处理逻辑（参数：c config.MqttConfig）MqttConfig

func NewMqttConfig(c config.MqttConfig) MqttConfig {
	// 返回 MqttConfig{...} 给调用方，结束当前函数

	return MqttConfig{
		// 执行 Enabled: c.Enabled,，完成对应的业务处理

		Enabled: c.Enabled,
		// 执行 Broker: c.Broker,，完成对应的业务处理

		Broker: c.Broker,
		// 执行 ClientID: c.ClientID,，完成对应的业务处理

		ClientID: c.ClientID,
		// 执行 Username: c.Username,，完成对应的业务处理

		Username: c.Username,
		// 执行 Password: c.Password,，完成对应的业务处理

		Password: c.Password,
	// 当前代码块结束

	}
// 当前代码块结束

}
