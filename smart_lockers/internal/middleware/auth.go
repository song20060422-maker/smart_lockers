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

type Claims struct {
	UserID   uint `json:"uid"`
	TokenVer int  `json:"tv"`
	jwt.RegisteredClaims
}

// GenerateToken 生成令牌 @ticket 06-个人中心
func GenerateToken(userID uint, tokenVer int) (string, error) {
	cfg := config.Get()
	now := time.Now()
	claims := Claims{
		UserID: userID, TokenVer: tokenVer,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析并校验令牌签名与有效期
func ParseToken(tokenStr string) (*Claims, error) {
	cfg := config.Get()
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.Secret), nil
	})
	if err != nil || !tok.Valid {
		return nil, err
	}
	return claims, nil
}

// Auth 鉴权中间件：从 Authorization: Bearer <token> 解析 uid/tv 注入上下文
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			pkg.AbortWithCode(c, pkg.ErrAuth)
			return
		}
		claims, err := ParseToken(strings.TrimPrefix(h, "Bearer "))
		if err != nil {
			pkg.AbortWithCode(c, pkg.ErrAuth)
			return
		}
		c.Set("uid", claims.UserID)
		c.Set("tv", claims.TokenVer)
		c.Next()
	}
}

// UID 从上下文取当前用户ID
func UID(c *gin.Context) uint {
	if v, ok := c.Get("uid"); ok {
		if u, ok := v.(uint); ok {
			return u
		}
	}
	return 0
}

// TV 从上下文取令牌版本
func TV(c *gin.Context) int {
	if v, ok := c.Get("tv"); ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}

// UIDStr 兼容字符串形式
func UIDStr(c *gin.Context) string { return strconv.FormatUint(uint64(UID(c)), 10) }
