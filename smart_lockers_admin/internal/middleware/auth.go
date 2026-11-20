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

type Claims struct {
	UserID uint              `json:"uid"`
	Role   model.AdminRole   `json:"role"`
	SiteID uint              `json:"sid"`
	jwt.RegisteredClaims
}

// GenerateToken 生成带角色的令牌 @ticket B端-RBAC
func GenerateToken(userID uint, role model.AdminRole, siteID uint) (string, error) {
	cfg := config.Get()
	now := time.Now()
	claims := Claims{
		UserID: userID, Role: role, SiteID: siteID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析校验令牌
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

// Auth 鉴权中间件：解析 Bearer token，注入 uid/role/site_id
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
		c.Set("role", claims.Role)
		c.Set("site_id", claims.SiteID)
		c.Next()
	}
}

// CurrentUser 当前登录身份
type CurrentUser struct {
	ID     uint
	Role   model.AdminRole
	SiteID uint
}

func User(c *gin.Context) CurrentUser {
	role, _ := c.Get("role")
	siteID, _ := c.Get("site_id")
	uid, _ := c.Get("uid")
	r, _ := role.(model.AdminRole)
	s, _ := siteID.(uint)
	u, _ := uid.(uint)
	return CurrentUser{ID: u, Role: r, SiteID: s}
}
