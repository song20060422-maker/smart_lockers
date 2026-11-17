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

var Rdb *redis.Client

func InitRedis(addr, password string, db int) {
	r := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := r.Ping(ctx).Err(); err != nil {
		log.Printf("[WARN] redis connect failed: %v (幂等/频控降级为本地内存实现)", err)
		Rdb = nil
		return
	}
	Rdb = r
	log.Println("[INFO] redis connected")
}

// ---------- 通用 KV 抽象（Redis 优先，降级内存），用于幂等键与频控 ----------

type kvEntry struct {
	value string
	exp   time.Time
}

var (
	memMu sync.Mutex
	memKV = make(map[string]kvEntry)
)

// SetNX 仅当键不存在时写入，成功返回 true（分布式锁/幂等基础）
func SetNX(key, value string, ttl time.Duration) bool {
	if Rdb != nil {
		ok, err := Rdb.SetNX(context.Background(), key, value, ttl).Result()
		if err == nil {
			return ok
		}
	}
	memMu.Lock()
	defer memMu.Unlock()
	if e, ok := memKV[key]; ok && time.Now().Before(e.exp) {
		return false
	}
	memKV[key] = kvEntry{value: value, exp: time.Now().Add(ttl)}
	return true
}

// Get 读取键值
func Get(key string) (string, bool) {
	if Rdb != nil {
		v, err := Rdb.Get(context.Background(), key).Result()
		if err == nil {
			return v, true
		}
	}
	memMu.Lock()
	defer memMu.Unlock()
	e, ok := memKV[key]
	if !ok || time.Now().After(e.exp) {
		return "", false
	}
	return e.value, true
}

// IncrWithTTL 计数并按窗口设置过期，用于频控
func IncrWithTTL(key string, ttl time.Duration) int64 {
	if Rdb != nil {
		pipe := Rdb.Pipeline()
		incr := pipe.Incr(context.Background(), key)
		pipe.Expire(context.Background(), key, ttl)
		_, _ = pipe.Exec(context.Background())
		return incr.Val()
	}
	memMu.Lock()
	defer memMu.Unlock()
	e, ok := memKV[key]
	if !ok || time.Now().After(e.exp) {
		memKV[key] = kvEntry{value: "1", exp: time.Now().Add(ttl)}
		return 1
	}
	var n int64
	fmt.Sscan(e.value, &n)
	n++
	e.value = fmt.Sprintf("%d", n)
	memKV[key] = e
	return n
}
