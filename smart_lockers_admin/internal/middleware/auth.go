// @ticket B端-RBAC：JWT 鉴权，Claims 携带 role/site_id，注入上下文
package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"smart_lockers_admin/internal/config"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
)

// 定义结构体 Claims：用于承载一组相关的字段数据

type Claims struct {
	// 结构体字段 UserID，类型为 uint，附带标签 `json:"uid"`

	UserID uint              `json:"uid"`
	// 结构体字段 Role，类型为 model.AdminRole，附带标签 `json:"role"`

	Role   model.AdminRole   `json:"role"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `json:"sid"`

	SiteID uint              `json:"sid"`
	// 执行 jwt.RegisteredClaims，完成对应的业务处理

	jwt.RegisteredClaims
// 当前代码块结束

}

// GenerateToken 生成带角色的令牌 @ticket B端-RBAC
// 定义函数 GenerateToken：封装一段可复用的处理逻辑（参数：userID uint, role model.AdminRole, siteID uint）

func GenerateToken(userID uint, role model.AdminRole, siteID uint) (string, error) {
	// 定义并初始化局部变量 cfg，取值为 config.Get()

	cfg := config.Get()
	// 定义并初始化局部变量 now，取值为 time.Now()

	now := time.Now()
	// 定义并初始化局部变量 claims，取值为 Claims{

	claims := Claims{
		// 执行 UserID: userID, Role: role, SiteID: siteID,，完成对应的业务处理

		UserID: userID, Role: role, SiteID: siteID,
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

// ParseToken 解析校验令牌
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

// Auth 鉴权中间件：解析 Bearer token，注入 uid/role/site_id
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
		// 执行 c.Set("role", claims.Role)，完成对应的业务处理

		c.Set("role", claims.Role)
		// 执行 c.Set("site_id", claims.SiteID)，完成对应的业务处理

		c.Set("site_id", claims.SiteID)
		// 执行 c.Next()，完成对应的业务处理

		c.Next()
	// 当前代码块结束

	}
// 当前代码块结束

}

// CurrentUser 当前登录身份
// 定义结构体 CurrentUser：用于承载一组相关的字段数据

type CurrentUser struct {
	// 结构体字段 ID，类型为 uint

	ID     uint
	// 结构体字段 Role，类型为 model.AdminRole

	Role   model.AdminRole
	// 结构体字段 SiteID，类型为 uint

	SiteID uint
// 当前代码块结束

}

// 定义函数 User：封装一段可复用的处理逻辑（参数：c *gin.Context）

func User(c *gin.Context) CurrentUser {
	// 定义并初始化局部变量 role, _，取值为 c.Get("role")

	role, _ := c.Get("role")
	// 定义并初始化局部变量 siteID, _，取值为 c.Get("site_id")

	siteID, _ := c.Get("site_id")
	// 定义并初始化局部变量 uid, _，取值为 c.Get("uid")

	uid, _ := c.Get("uid")
	// 定义并初始化局部变量 r, _，取值为 role.(model.AdminRole)

	r, _ := role.(model.AdminRole)
	// 定义并初始化局部变量 s, _，取值为 siteID.(uint)

	s, _ := siteID.(uint)
	// 定义并初始化局部变量 u, _，取值为 uid.(uint)

	u, _ := uid.(uint)
	// 返回 CurrentUser{ID: u, Role: r, SiteID: s} 给调用方，结束当前函数

	return CurrentUser{ID: u, Role: r, SiteID: s}
// 当前代码块结束

}
