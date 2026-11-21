// @ticket 00-project：Redis 初始化；不可用降级为内存实现，保证 go run 无 Redis 也能启动
package pkg

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// 声明变量 Rdb *redis.Client，用于保存运行过程中的数据

var Rdb *redis.Client

// 定义函数 InitRedis：封装一段可复用的处理逻辑（参数：addr, password string, db int）

func InitRedis(addr, password string, db int) {
	// 定义并初始化局部变量 r，取值为 redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})

	r := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	// 定义并初始化局部变量 ctx, cancel，取值为 context.WithTimeout(context.Background(), 2*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// 延迟注册 cancel()，在当前函数返回前自动执行

	defer cancel()
	// 条件判断：err := r.Ping(ctx).Err(); err != nil 时执行下面逻辑

	if err := r.Ping(ctx).Err(); err != nil {
		// 执行 log.Printf("[WARN] redis connect failed: %v (幂等/频控降级为本地内存实现)", err)，完成对应的业务处理

		log.Printf("[WARN] redis connect failed: %v (幂等/频控降级为本地内存实现)", err)
		// 为 Rdb 赋值：nil

		Rdb = nil
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 为 Rdb 赋值：r

	Rdb = r
	// 执行 log.Println("[INFO] redis connected")，完成对应的业务处理

	log.Println("[INFO] redis connected")
// 当前代码块结束

}

// ---------- 通用 KV 抽象（Redis 优先，降级内存），用于幂等键与频控 ----------

// 定义结构体 kvEntry：用于承载一组相关的字段数据

type kvEntry struct {
	// 结构体字段 value，类型为 string

	value string
	// 结构体字段 exp，类型为 time.Time

	exp   time.Time
// 当前代码块结束

}

// 开始定义一组包级/局部变量的 var 声明块

var (
	// 执行 memMu sync.Mutex，完成对应的业务处理

	memMu sync.Mutex
	// 为 memKV 赋值：make(map[string]kvEntry)

	memKV = make(map[string]kvEntry)
// 结束参数/表达式列表

)

// SetNX 仅当键不存在时写入，成功返回 true（分布式锁/幂等基础）
// 定义函数 SetNX：封装一段可复用的处理逻辑（参数：key, value string, ttl time.Duration）

func SetNX(key, value string, ttl time.Duration) bool {
	// 条件判断：Rdb != nil 时执行下面逻辑

	if Rdb != nil {
		// 定义并初始化局部变量 ok, err，取值为 Rdb.SetNX(context.Background(), key, value, ttl).Result()

		ok, err := Rdb.SetNX(context.Background(), key, value, ttl).Result()
		// 条件判断：err == nil 时执行下面逻辑

		if err == nil {
			// 返回 ok 给调用方，结束当前函数

			return ok
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 执行 memMu.Lock()，完成对应的业务处理

	memMu.Lock()
	// 延迟注册 memMu.Unlock()，在当前函数返回前自动执行

	defer memMu.Unlock()
	// 条件判断：e, ok := memKV[key]; ok && time.Now().Before(e.exp) 时执行下面逻辑

	if e, ok := memKV[key]; ok && time.Now().Before(e.exp) {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 为 memKV[key] 赋值：kvEntry{value: value, exp: time.Now().Add(ttl)}

	memKV[key] = kvEntry{value: value, exp: time.Now().Add(ttl)}
	// 返回 true 给调用方，结束当前函数

	return true
// 当前代码块结束

}

// Get 读取键值
// 定义函数 Get：封装一段可复用的处理逻辑（参数：key string）

func Get(key string) (string, bool) {
	// 条件判断：Rdb != nil 时执行下面逻辑

	if Rdb != nil {
		// 定义并初始化局部变量 v, err，取值为 Rdb.Get(context.Background(), key).Result()

		v, err := Rdb.Get(context.Background(), key).Result()
		// 条件判断：err == nil 时执行下面逻辑

		if err == nil {
			// 返回 v, true 给调用方，结束当前函数

			return v, true
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 执行 memMu.Lock()，完成对应的业务处理

	memMu.Lock()
	// 延迟注册 memMu.Unlock()，在当前函数返回前自动执行

	defer memMu.Unlock()
	// 定义并初始化局部变量 e, ok，取值为 memKV[key]

	e, ok := memKV[key]
	// 条件判断：!ok || time.Now().After(e.exp) 时执行下面逻辑

	if !ok || time.Now().After(e.exp) {
		// 返回 "", false 给调用方，结束当前函数

		return "", false
	// 当前代码块结束

	}
	// 返回 e.value, true 给调用方，结束当前函数

	return e.value, true
// 当前代码块结束

}

// IncrWithTTL 计数并按窗口设置过期，用于频控
// 定义函数 IncrWithTTL：封装一段可复用的处理逻辑（参数：key string, ttl time.Duration）

func IncrWithTTL(key string, ttl time.Duration) int64 {
	// 条件判断：Rdb != nil 时执行下面逻辑

	if Rdb != nil {
		// 定义并初始化局部变量 pipe，取值为 Rdb.Pipeline()

		pipe := Rdb.Pipeline()
		// 定义并初始化局部变量 incr，取值为 pipe.Incr(context.Background(), key)

		incr := pipe.Incr(context.Background(), key)
		// 执行 pipe.Expire(context.Background(), key, ttl)，完成对应的业务处理

		pipe.Expire(context.Background(), key, ttl)
		// 为 _, _ 赋值：pipe.Exec(context.Background())

		_, _ = pipe.Exec(context.Background())
		// 返回 incr.Val() 给调用方，结束当前函数

		return incr.Val()
	// 当前代码块结束

	}
	// 执行 memMu.Lock()，完成对应的业务处理

	memMu.Lock()
	// 延迟注册 memMu.Unlock()，在当前函数返回前自动执行

	defer memMu.Unlock()
	// 定义并初始化局部变量 e, ok，取值为 memKV[key]

	e, ok := memKV[key]
	// 条件判断：!ok || time.Now().After(e.exp) 时执行下面逻辑

	if !ok || time.Now().After(e.exp) {
		// 为 memKV[key] 赋值：kvEntry{value: "1", exp: time.Now().Add(ttl)}

		memKV[key] = kvEntry{value: "1", exp: time.Now().Add(ttl)}
		// 返回 1 给调用方，结束当前函数

		return 1
	// 当前代码块结束

	}
	// 声明变量 n int64，用于保存运行过程中的数据

	var n int64
	// 执行 fmt.Sscan(e.value, &n)，完成对应的业务处理

	fmt.Sscan(e.value, &n)
	// 执行 n++，完成对应的业务处理

	n++
	// 为 e.value 赋值：fmt.Sprintf("%d", n)

	e.value = fmt.Sprintf("%d", n)
	// 为 memKV[key] 赋值：e

	memKV[key] = e
	// 返回 n 给调用方，结束当前函数

	return n
// 当前代码块结束

}
