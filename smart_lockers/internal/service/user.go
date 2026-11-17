// @ticket 06-个人中心：头像上传校验、我的服务动态配置、卡券校验、退出登录即时失效
package service

import (
	"strings"
	"time"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
)

// allowedAvatarTypes 头像允许格式 @ticket 06-个人中心
var allowedAvatarTypes = map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true}

// ValidateAvatar 头像格式/大小校验 @ticket 06-个人中心 P95≤500ms
func ValidateAvatar(contentType string, sizeBytes int64) int {
	if !allowedAvatarTypes[strings.ToLower(contentType)] {
		return pkg.ErrAvatarType
	}
	maxBytes := int64(5 * 1024 * 1024) // 5MB
	if sizeBytes > maxBytes || sizeBytes <= 0 {
		return pkg.ErrAvatarSize
	}
	return pkg.CodeOK
}

// UpdateAvatar 更新头像（校验通过后写库）@ticket 06-个人中心
func UpdateAvatar(userID uint, contentType string, sizeBytes int64, url string) (map[string]interface{}, int) {
	if code := ValidateAvatar(contentType, sizeBytes); code != pkg.CodeOK {
		return nil, code
	}
	if pkg.DB != nil {
		if err := pkg.DB.Model(&model.UserProfile{}).Where("id = ?", userID).Update("avatar", url).Error; err != nil {
			return nil, pkg.ErrInternal
		}
	}
	return map[string]interface{}{"avatar": url}, pkg.CodeOK
}

// MyServices 我的服务：按角色动态配置 @ticket 06-个人中心
func MyServices(role string) ([]model.MyService, int) {
	var all []model.MyService
	if pkg.DB == nil {
		return []model.MyService{}, pkg.CodeOK
	}
	pkg.DB.Where("enabled = ?", true).Find(&all)
	res := make([]model.MyService, 0, len(all))
	for _, s := range all {
		// roles 为空表示全部角色可见；否则需命中角色 @ticket 06-个人中心
		if s.Roles == "" || strings.Contains(s.Roles, role) {
			res = append(res, s)
		}
	}
	return res, pkg.CodeOK
}

// Coupons 卡券列表（可用/过期实时校验）@ticket 06-个人中心
func Coupons(userID uint) (map[string]interface{}, int) {
	available := make([]model.Coupon, 0)
	expired := make([]model.Coupon, 0)
	if pkg.DB != nil {
		var list []model.Coupon
		pkg.DB.Where("user_id = ?", userID).Find(&list)
		now := time.Now()
		for _, c := range list {
			if c.Used || now.After(c.ExpireAt) {
				expired = append(expired, c)
			} else {
				available = append(available, c)
			}
		}
	}
	return map[string]interface{}{"available": available, "expired": expired}, pkg.CodeOK
}

// ValidateCoupon 下单前的卡券实时校验（幂等抵扣由订单唯一键保证）@ticket 06-个人中心
func ValidateCoupon(code string, userID uint) (int64, int) {
	if pkg.DB == nil {
		return 0, pkg.ErrCouponExp
	}
	var c model.Coupon
	if err := pkg.DB.Where("code = ? AND user_id = ?", code, userID).First(&c).Error; err != nil {
		return 0, pkg.ErrCouponExp
	}
	if c.Used || time.Now().After(c.ExpireAt) {
		return 0, pkg.ErrCouponExp
	}
	return c.Discount, pkg.CodeOK
}

// ConsumeCoupon 幂等抵扣卡券 @ticket 06-个人中心 + 04-物品寄存
func ConsumeCoupon(code string, userID uint, bizNo string) int {
	if pkg.DB == nil {
		return pkg.ErrCouponExp
	}
	res := pkg.DB.Model(&model.Coupon{}).
		Where("code = ? AND user_id = ? AND used = 0 AND expire_at > ?", code, userID, time.Now()).
		Update("used", true)
	if res.Error != nil || res.RowsAffected == 0 {
		return pkg.ErrCouponExp
	}
	pkg.Logf("[AUDIT] @ticket 06 coupon consumed: code=%s uid=%d biz=%s", code, userID, bizNo)
	return pkg.CodeOK
}

// Logout 退出登录：tokenVer 自增使旧 token 立即失效 + 幂等 @ticket 06-个人中心
func Logout(userID uint, currentVer int) (map[string]interface{}, int) {
	if pkg.DB != nil {
		res := pkg.DB.Model(&model.UserProfile{}).Where("id = ? AND token_ver = ?", userID, currentVer).
			Update("token_ver", currentVer+1)
		if res.Error != nil && res.Error.Error() != "record not found" {
			return nil, pkg.ErrInternal
		}
	}
	return map[string]interface{}{"logout": true}, pkg.CodeOK
}

// IssueToken 登录签发令牌（读取用户当前 tokenVer）@ticket 06-个人中心
func IssueToken(userID uint) (string, int) {
	ver := 0
	if pkg.DB != nil {
		var u model.UserProfile
		if pkg.DB.First(&u, userID).Error == nil {
			ver = u.TokenVer
		}
	}
	tok, err := middleware.GenerateToken(userID, ver)
	if err != nil {
		return "", pkg.ErrInternal
	}
	return tok, pkg.CodeOK
}
