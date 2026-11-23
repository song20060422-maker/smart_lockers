// @ticket B端-IoT中台：EMQX 轻客户端（发布 platform/# 指令 + 订阅 device/response、device/alert 记录回执）
// 对齐方案 3.4 ACL：后端服务允许发布 platform/#、订阅 device/#；不可达时降级本地演示模式
package pkg

import (
	"fmt"
	"log"
	"sync"
	"time"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"

	"smart_lockers_admin/internal/config"
)

// MqttConfig B端 MQTT 连接参数（从 config.MqttConfig 转换）
// 定义结构体 MqttConfig：用于承载一组相关的字段数据

type MqttConfig struct {
	// 结构体字段 Enabled，类型为 bool

	Enabled bool
	// 结构体字段 Broker，类型为 string

	Broker string
	// 结构体字段 ClientID，类型为 string

	ClientID string
	// 结构体字段 Username，类型为 string

	Username string
	// 结构体字段 Password，类型为 string

	Password string
// 当前代码块结束

}

// 声明变量 mqttMu/mqttClient/mqttReady，用于保存运行过程中的数据

var (
	// 执行 mqttMu sync.Mutex，完成对应的业务处理

	mqttMu sync.Mutex
	// 执行 mqttClient pahomqtt.Client，完成对应的业务处理

	mqttClient pahomqtt.Client
	// 执行 mqttReady bool，完成对应的业务处理

	mqttReady bool
// 结束参数/表达式列表

)

// InitMQTT 初始化 B端 MQTT 客户端；handler 为上行回执/告警统一入口（topic,payload）
// 定义函数 InitMQTT：封装一段可复用的处理逻辑（参数：cfg MqttConfig, handler func(topic string, payload []byte)）

func InitMQTT(cfg MqttConfig, handler func(topic string, payload []byte)) {
	// 加互斥锁保护初始化
	// 执行 mqttMu.Lock()，完成对应的业务处理

	mqttMu.Lock()
	// 延迟注册解锁

	defer mqttMu.Unlock()
	// 条件判断：!cfg.Enabled 时执行下面逻辑

	if !cfg.Enabled {
		// 执行 log.Println("[INFO] admin mqtt disabled (本地演示模式，指令仅落库不发布)")，完成对应的业务处理

		log.Println("[INFO] admin mqtt disabled (本地演示模式，指令仅落库不发布)")
		// 为 mqttReady 赋值：false

		mqttReady = false
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 构建连接参数
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
	// 自动重连 + 保留会话（离线回执不丢失）
	// 执行 opts.SetAutoReconnect(true)，完成对应的业务处理

	opts.SetAutoReconnect(true)
	// 执行 opts.SetCleanSession(false)，完成对应的业务处理

	opts.SetCleanSession(false)
	// 执行 opts.SetKeepAlive(60 * time.Second)，完成对应的业务处理

	opts.SetKeepAlive(60 * time.Second)
	// 连接成功回调：标记就绪
	// 执行 opts.SetOnConnectHandler(func(c pahomqtt.Client) { mqttReady = true })，完成对应的业务处理

	opts.SetOnConnectHandler(func(c pahomqtt.Client) {
		// 加锁置位
		// 执行 mqttMu.Lock()，完成对应的业务处理

		mqttMu.Lock()
		// 为 mqttReady 赋值：true

		mqttReady = true
		// 解锁

		mqttMu.Unlock()
		// 执行 log.Println("[INFO] admin mqtt connected")，完成对应的业务处理

		log.Println("[INFO] admin mqtt connected")
	// 当前代码块结束

	})
	// 连接丢失回调：记录日志
	// 执行 opts.SetConnectionLostHandler(...)，完成对应的业务处理

	opts.SetConnectionLostHandler(func(c pahomqtt.Client, err error) {
		// 执行 log.Printf("[WARN] admin mqtt connection lost: %v", err)，完成对应的业务处理

		log.Printf("[WARN] admin mqtt connection lost: %v", err)
	// 当前代码块结束

	})
	// 设置统一上行消息入口
	// 执行 opts.SetDefaultPublishHandler(func(c pahomqtt.Client, m pahomqtt.Message) { handler(m.Topic(), m.Payload()) })，完成对应的业务处理

	opts.SetDefaultPublishHandler(func(c pahomqtt.Client, m pahomqtt.Message) {
		// 转发到业务回调
		// 执行 handler(m.Topic(), m.Payload())，完成对应的业务处理

		handler(m.Topic(), m.Payload())
	// 当前代码块结束

	})
	// 建立客户端并异步连接
	// 定义并初始化局部变量 cli，取值为 pahomqtt.NewClient(opts)

	cli := pahomqtt.NewClient(opts)
	// 定义并初始化局部变量 token，取值为 cli.Connect()

	token := cli.Connect()
	// 条件判断：!token.WaitTimeout(3*time.Second) || token.Error() != nil 时执行下面逻辑

	if ok := token.WaitTimeout(3 * time.Second); !ok || token.Error() != nil {
		// 执行 log.Printf("[WARN] admin mqtt connect failed: %v (降级为本地演示模式)", token.Error())，完成对应的业务处理

		log.Printf("[WARN] admin mqtt connect failed: %v (降级为本地演示模式)", token.Error())
		// 为 mqttReady 赋值：false

		mqttReady = false
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 为 mqttClient 赋值：cli

	mqttClient = cli
	// 执行 log.Printf("[INFO] admin mqtt client ready (broker=%s)", cfg.Broker)，完成对应的业务处理

	log.Printf("[INFO] admin mqtt client ready (broker=%s)", cfg.Broker)
// 当前代码块结束

}

// MqttReady 返回 B端 MQTT 是否可用
// 定义函数 MqttReady：封装一段可复用的处理逻辑（参数：）bool

func MqttReady() bool {
	// 加锁读取
	// 执行 mqttMu.Lock()，完成对应的业务处理

	mqttMu.Lock()
	// 延迟注册解锁

	defer mqttMu.Unlock()
	// 返回 mqttReady 给调用方，结束当前函数

	return mqttReady
// 当前代码块结束

}

// MqttPublish 发布下行消息（platform/# 指令/OTA）
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
		// 返回发布错误给调用方，结束当前函数

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

// MqttSubscribe 订阅上行 Topic（如 device/response/#、device/alert/#）
// 定义函数 MqttSubscribe：封装一段可复用的处理逻辑（参数：topic string）error

func MqttSubscribe(topic string) error {
	// 条件判断：!MqttReady() 时执行下面逻辑

	if !MqttReady() {
		// 返回 fmt.Errorf("mqtt unavailable") 给调用方，结束当前函数

		return fmt.Errorf("mqtt unavailable")
	// 当前代码块结束

	}
	// 订阅（nil 回调走 DefaultPublishHandler 统一入口）
	// 定义并初始化局部变量 t，取值为 mqttClient.Subscribe(topic, 1, nil)

	t := mqttClient.Subscribe(topic, 1, nil)
	// 条件判断：!t.WaitTimeout(3*time.Second) || t.Error() != nil 时执行下面逻辑

	if ok := t.WaitTimeout(3 * time.Second); !ok || t.Error() != nil {
		// 返回订阅错误给调用方，结束当前函数

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

// BrokerURL 返回配置的 Broker 地址（总览展示用）
// 定义函数 BrokerURL：封装一段可复用的处理逻辑（参数：）string

func BrokerURL() string {
	// 返回配置 Broker 给调用方，结束当前函数

	return config.Get().MQTT.Broker
// 当前代码块结束

}

// NewMqttConfig 从 B端全局配置构造 MqttConfig
// 定义函数 NewMqttConfig：封装一段可复用的处理逻辑（参数：c config.MqttConfig）MqttConfig

func NewMqttConfig(c config.MqttConfig) MqttConfig {
	// 返回转换结果给调用方，结束当前函数

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
