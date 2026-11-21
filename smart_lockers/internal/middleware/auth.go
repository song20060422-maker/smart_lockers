// @ticket 00-project + 06-个人中心：JWT 鉴权；退出登录通过 tokenVer 校验使旧令牌立即失效
package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"smart_lockers/internal/config"
	"smart_lockers/internal/pkg"
)

// 定义结构体 Claims：用于承载一组相关的字段数据

type Claims struct {
	// 结构体字段 UserID，类型为 uint，附带标签 `json:"uid"`

	UserID   uint `json:"uid"`
	// 结构体字段 TokenVer，类型为 int，附带标签 `json:"tv"`

	TokenVer int  `json:"tv"`
	// 执行 jwt.RegisteredClaims，完成对应的业务处理

	jwt.RegisteredClaims
// 当前代码块结束

}

// GenerateToken 生成令牌 @ticket 06-个人中心
// 定义函数 GenerateToken：封装一段可复用的处理逻辑（参数：userID uint, tokenVer int）

func GenerateToken(userID uint, tokenVer int) (string, error) {
	// 定义并初始化局部变量 cfg，取值为 config.Get()

	cfg := config.Get()
	// 定义并初始化局部变量 now，取值为 time.Now()

	now := time.Now()
	// 定义并初始化局部变量 claims，取值为 Claims{

	claims := Claims{
		// 执行 UserID: userID, TokenVer: tokenVer,，完成对应的业务处理

		UserID: userID, TokenVer: tokenVer,
		// 执行 RegisteredClaims: jwt.RegisteredClaims{，完成对应的业务处理

		RegisteredClaims: jwt.RegisteredClaims{
			// 执行 ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.ExpireHours) * time.Hour)),，完成对应的业务处理

			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.ExpireHours) * time.Hour)),
			// 执行 IssuedAt:  jwt.NewNumericDate(now),，完成对应的业务处理

			IssuedAt:  jwt.NewNumericDate(now),
		// 定义结束

		},
	// 当前代码块结束

	}
	// 返回 jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWT.Secret)) 给调用方，结束当前函数

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWT.Secret))
// 当前代码块结束

}

// ParseToken 解析并校验令牌签名与有效期
// 定义函数 ParseToken：封装一段可复用的处理逻辑（参数：tokenStr string）

func ParseToken(tokenStr string) (*Claims, error) {
	// 定义并初始化局部变量 cfg，取值为 config.Get()

	cfg := config.Get()
	// 定义并初始化局部变量 claims，取值为 &Claims{}

	claims := &Claims{}
	// 定义并初始化局部变量 tok, err，取值为 jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {

	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		// 返回 []byte(cfg.JWT.Secret), nil 给调用方，结束当前函数

		return []byte(cfg.JWT.Secret), nil
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err != nil || !tok.Valid 时执行下面逻辑

	if err != nil || !tok.Valid {
		// 返回 nil, err 给调用方，结束当前函数

		return nil, err
	// 当前代码块结束

	}
	// 返回 claims, nil 给调用方，结束当前函数

	return claims, nil
// 当前代码块结束

}

// Auth 鉴权中间件：从 Authorization: Bearer <token> 解析 uid/tv 注入上下文
// 定义函数 Auth：封装一段可复用的处理逻辑（参数：）

func Auth() gin.HandlerFunc {
	// 返回 func(c *gin.Context) { 给调用方，结束当前函数

	return func(c *gin.Context) {
		// 定义并初始化局部变量 h，取值为 c.GetHeader("Authorization")

		h := c.GetHeader("Authorization")
		// 条件判断：!strings.HasPrefix(h, "Bearer ") 时执行下面逻辑

		if !strings.HasPrefix(h, "Bearer ") {
			// 执行 pkg.AbortWithCode(c, pkg.ErrAuth)，完成对应的业务处理

			pkg.AbortWithCode(c, pkg.ErrAuth)
			// 返回  给调用方，结束当前函数

			return
		// 当前代码块结束

		}
		// 定义并初始化局部变量 claims, err，取值为 ParseToken(strings.TrimPrefix(h, "Bearer "))

		claims, err := ParseToken(strings.TrimPrefix(h, "Bearer "))
		// 条件判断：err != nil 时执行下面逻辑

		if err != nil {
			// 执行 pkg.AbortWithCode(c, pkg.ErrAuth)，完成对应的业务处理

			pkg.AbortWithCode(c, pkg.ErrAuth)
			// 返回  给调用方，结束当前函数

			return
		// 当前代码块结束

		}
		// 执行 c.Set("uid", claims.UserID)，完成对应的业务处理

		c.Set("uid", claims.UserID)
		// 执行 c.Set("tv", claims.TokenVer)，完成对应的业务处理

		c.Set("tv", claims.TokenVer)
		// 执行 c.Next()，完成对应的业务处理

		c.Next()
	// 当前代码块结束

	}
// 当前代码块结束

}

// UID 从上下文取当前用户ID
// 定义函数 UID：封装一段可复用的处理逻辑（参数：c *gin.Context）

func UID(c *gin.Context) uint {
	// 条件判断：v, ok := c.Get("uid"); ok 时执行下面逻辑

	if v, ok := c.Get("uid"); ok {
		// 条件判断：u, ok := v.(uint); ok 时执行下面逻辑

		if u, ok := v.(uint); ok {
			// 返回 u 给调用方，结束当前函数

			return u
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 0 给调用方，结束当前函数

	return 0
// 当前代码块结束

}

// TV 从上下文取令牌版本
// 定义函数 TV：封装一段可复用的处理逻辑（参数：c *gin.Context）

func TV(c *gin.Context) int {
	// 条件判断：v, ok := c.Get("tv"); ok 时执行下面逻辑

	if v, ok := c.Get("tv"); ok {
		// 条件判断：i, ok := v.(int); ok 时执行下面逻辑

		if i, ok := v.(int); ok {
			// 返回 i 给调用方，结束当前函数

			return i
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 0 给调用方，结束当前函数

	return 0
// 当前代码块结束

}

// UIDStr 兼容字符串形式
// 定义函数 UIDStr：封装一段可复用的处理逻辑（参数：c *gin.Context）

func UIDStr(c *gin.Context) string { return strconv.FormatUint(uint64(UID(c)), 10) }
