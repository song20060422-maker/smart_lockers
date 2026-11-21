// @ticket 06-个人中心：头像上传校验、我的服务动态配置、卡券校验、退出登录即时失效
package service

import (
	"regexp"
	"strings"
	"time"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
)

// allowedAvatarTypes 头像允许格式 @ticket 06-个人中心
// 声明变量 allowedAvatarTypes = map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true}，用于保存运行过程中的数据

var allowedAvatarTypes = map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true}

// ValidateAvatar 头像格式/大小校验 @ticket 06-个人中心 P95≤500ms
// 定义函数 ValidateAvatar：封装一段可复用的处理逻辑（参数：contentType string, sizeBytes int64）

func ValidateAvatar(contentType string, sizeBytes int64) int {
	// 条件判断：!allowedAvatarTypes[strings.ToLower(contentType)] 时执行下面逻辑

	if !allowedAvatarTypes[strings.ToLower(contentType)] {
		// 返回 pkg.ErrAvatarType 给调用方，结束当前函数

		return pkg.ErrAvatarType
	// 当前代码块结束

	}
	// 定义并初始化局部变量 maxBytes，取值为 int64(5 * 1024 * 1024) // 5MB

	maxBytes := int64(5 * 1024 * 1024) // 5MB
	// 条件判断：sizeBytes > maxBytes || sizeBytes <= 0 时执行下面逻辑

	if sizeBytes > maxBytes || sizeBytes <= 0 {
		// 返回 pkg.ErrAvatarSize 给调用方，结束当前函数

		return pkg.ErrAvatarSize
	// 当前代码块结束

	}
	// 返回 pkg.CodeOK 给调用方，结束当前函数

	return pkg.CodeOK
// 当前代码块结束

}

// UpdateAvatar 更新头像（校验通过后写库）@ticket 06-个人中心
// 定义函数 UpdateAvatar：封装一段可复用的处理逻辑（参数：userID uint, contentType string, sizeBytes int64, url string）

func UpdateAvatar(userID uint, contentType string, sizeBytes int64, url string) (map[string]interface{}, int) {
	// 条件判断：code := ValidateAvatar(contentType, sizeBytes); code != pkg.CodeOK 时执行下面逻辑

	if code := ValidateAvatar(contentType, sizeBytes); code != pkg.CodeOK {
		// 返回 nil, code 给调用方，结束当前函数

		return nil, code
	// 当前代码块结束

	}
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 条件判断：err := pkg.DB.Model(&model.UserProfile{}).Where("id = ?", userID).Update("avatar", url).Error; err != nil 时执行下面逻辑

		if err := pkg.DB.Model(&model.UserProfile{}).Where("id = ?", userID).Update("avatar", url).Error; err != nil {
			// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

			return nil, pkg.ErrInternal
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"avatar": url}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"avatar": url}, pkg.CodeOK
// 当前代码块结束

}

// MyServices 我的服务：按角色动态配置 @ticket 06-个人中心
// 定义函数 MyServices：封装一段可复用的处理逻辑（参数：role string）

func MyServices(role string) ([]model.MyService, int) {
	// 声明变量 all []model.MyService，用于保存运行过程中的数据

	var all []model.MyService
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 []model.MyService{}, pkg.CodeOK 给调用方，结束当前函数

		return []model.MyService{}, pkg.CodeOK
	// 当前代码块结束

	}
	// 为 pkg.DB.Where("enabled 赋值：?", true).Find(&all)

	pkg.DB.Where("enabled = ?", true).Find(&all)
	// 定义并初始化局部变量 res，取值为 make([]model.MyService, 0, len(all))

	res := make([]model.MyService, 0, len(all))
	// 循环处理：for _, s := range all，逐轮迭代

	for _, s := range all {
		// roles 为空表示全部角色可见；否则需命中角色 @ticket 06-个人中心
		// 条件判断：s.Roles == "" || strings.Contains(s.Roles, role) 时执行下面逻辑

		if s.Roles == "" || strings.Contains(s.Roles, role) {
			// 为 res 赋值：append(res, s)

			res = append(res, s)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 res, pkg.CodeOK 给调用方，结束当前函数

	return res, pkg.CodeOK
// 当前代码块结束

}

// Coupons 卡券列表（可用/过期实时校验）@ticket 06-个人中心
// 定义函数 Coupons：封装一段可复用的处理逻辑（参数：userID uint）

func Coupons(userID uint) (map[string]interface{}, int) {
	// 定义并初始化局部变量 available，取值为 make([]model.Coupon, 0)

	available := make([]model.Coupon, 0)
	// 定义并初始化局部变量 expired，取值为 make([]model.Coupon, 0)

	expired := make([]model.Coupon, 0)
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 声明变量 list []model.Coupon，用于保存运行过程中的数据

		var list []model.Coupon
		// 为 pkg.DB.Where("user_id 赋值：?", userID).Find(&list)

		pkg.DB.Where("user_id = ?", userID).Find(&list)
		// 定义并初始化局部变量 now，取值为 time.Now()

		now := time.Now()
		// 循环处理：for _, c := range list，逐轮迭代

		for _, c := range list {
			// 条件判断：c.Used || now.After(c.ExpireAt) 时执行下面逻辑

			if c.Used || now.After(c.ExpireAt) {
				// 为 expired 赋值：append(expired, c)

				expired = append(expired, c)
			// 上述条件不成立时，执行 else 分支逻辑

			} else {
				// 为 available 赋值：append(available, c)

				available = append(available, c)
			// 当前代码块结束

			}
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"available": available, "expired": expired}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"available": available, "expired": expired}, pkg.CodeOK
// 当前代码块结束

}

// ValidateCoupon 下单前的卡券实时校验（幂等抵扣由订单唯一键保证）@ticket 06-个人中心
// 定义函数 ValidateCoupon：封装一段可复用的处理逻辑（参数：code string, userID uint）

func ValidateCoupon(code string, userID uint) (int64, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 0, pkg.ErrCouponExp 给调用方，结束当前函数

		return 0, pkg.ErrCouponExp
	// 当前代码块结束

	}
	// 声明变量 c model.Coupon，用于保存运行过程中的数据

	var c model.Coupon
	// 条件判断：err := pkg.DB.Where("code = ? AND user_id = ?", code, userID).First(&c).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("code = ? AND user_id = ?", code, userID).First(&c).Error; err != nil {
		// 返回 0, pkg.ErrCouponExp 给调用方，结束当前函数

		return 0, pkg.ErrCouponExp
	// 当前代码块结束

	}
	// 条件判断：c.Used || time.Now().After(c.ExpireAt) 时执行下面逻辑

	if c.Used || time.Now().After(c.ExpireAt) {
		// 返回 0, pkg.ErrCouponExp 给调用方，结束当前函数

		return 0, pkg.ErrCouponExp
	// 当前代码块结束

	}
	// 返回 c.Discount, pkg.CodeOK 给调用方，结束当前函数

	return c.Discount, pkg.CodeOK
// 当前代码块结束

}

// ConsumeCoupon 幂等抵扣卡券 @ticket 06-个人中心 + 04-物品寄存
// 定义函数 ConsumeCoupon：封装一段可复用的处理逻辑（参数：code string, userID uint, bizNo string）

func ConsumeCoupon(code string, userID uint, bizNo string) int {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 pkg.ErrCouponExp 给调用方，结束当前函数

		return pkg.ErrCouponExp
	// 当前代码块结束

	}
	// 定义并初始化局部变量 res，取值为 pkg.DB.Model(&model.Coupon{}).

	res := pkg.DB.Model(&model.Coupon{}).
		// 为 Where("code 赋值：? AND user_id = ? AND used = 0 AND expire_at > ?", code, userID, time.Now()).

		Where("code = ? AND user_id = ? AND used = 0 AND expire_at > ?", code, userID, time.Now()).
		// 执行 Update("used", true)，完成对应的业务处理

		Update("used", true)
	// 条件判断：res.Error != nil || res.RowsAffected == 0 时执行下面逻辑

	if res.Error != nil || res.RowsAffected == 0 {
		// 返回 pkg.ErrCouponExp 给调用方，结束当前函数

		return pkg.ErrCouponExp
	// 当前代码块结束

	}
	// 执行 pkg.Logf("[AUDIT] @ticket 06 coupon consumed: code=%s uid=%d biz=%s", code, userID, bizNo)，完成对应的业务处理

	pkg.Logf("[AUDIT] @ticket 06 coupon consumed: code=%s uid=%d biz=%s", code, userID, bizNo)
	// 返回 pkg.CodeOK 给调用方，结束当前函数

	return pkg.CodeOK
// 当前代码块结束

}

// Logout 退出登录：tokenVer 自增使旧 token 立即失效 + 幂等 @ticket 06-个人中心
// 定义函数 Logout：封装一段可复用的处理逻辑（参数：userID uint, currentVer int）

func Logout(userID uint, currentVer int) (map[string]interface{}, int) {
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 定义并初始化局部变量 res，取值为 pkg.DB.Model(&model.UserProfile{}).Where("id = ? AND token_ver = ?", userID, currentVer).

		res := pkg.DB.Model(&model.UserProfile{}).Where("id = ? AND token_ver = ?", userID, currentVer).
			// 执行 Update("token_ver", currentVer+1)，完成对应的业务处理

			Update("token_ver", currentVer+1)
		// 条件判断：res.Error != nil && res.Error.Error() != "record not found" 时执行下面逻辑

		if res.Error != nil && res.Error.Error() != "record not found" {
			// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

			return nil, pkg.ErrInternal
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"logout": true}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"logout": true}, pkg.CodeOK
// 当前代码块结束

}

// IssueToken 登录签发令牌（读取用户当前 tokenVer）@ticket 06-个人中心
// 定义函数 IssueToken：封装一段可复用的处理逻辑（参数：userID uint）

func IssueToken(userID uint) (string, int) {
	// 定义并初始化局部变量 ver，取值为 0

	ver := 0
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 声明变量 u model.UserProfile，用于保存运行过程中的数据

		var u model.UserProfile
		// 条件判断：pkg.DB.First(&u, userID).Error == nil 时执行下面逻辑

		if pkg.DB.First(&u, userID).Error == nil {
			// 为 ver 赋值：u.TokenVer

			ver = u.TokenVer
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 定义并初始化局部变量 tok, err，取值为 middleware.GenerateToken(userID, ver)

	tok, err := middleware.GenerateToken(userID, ver)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 "", pkg.ErrInternal 给调用方，结束当前函数

		return "", pkg.ErrInternal
	// 当前代码块结束

	}
	// 返回 tok, pkg.CodeOK 给调用方，结束当前函数

	return tok, pkg.CodeOK
// 当前代码块结束

}

// DemoLogin @ticket 00-project：可验证的演示登录——按手机号定位/创建用户并签发真实 JWT，替代前端本地生成 token
// 演示规则：手机号 11 位、验证码 4~6 位数字即通过（与前端提示一致），不依赖短信服务
// 定义函数 DemoLogin：封装一段可复用的处理逻辑（参数：phone string, code string）

func DemoLogin(phone string, code string) (uint, string, int) {
	// 校验手机号格式：1 开头的 11 位数字

	if !regexp.MustCompile(`^1\d{10}$`).MatchString(phone) {
		// 手机号不合法返回参数错误

		return 0, "", pkg.ErrBad
	// 当前代码块结束

	}
	// 校验验证码格式：4~6 位数字

	if !regexp.MustCompile(`^\d{4,6}$`).MatchString(code) {
		// 验证码不合法返回参数错误

		return 0, "", pkg.ErrBad
	// 当前代码块结束

	}
	// 默认演示 uid，与既有演示用户数据一致（10001）

	var uid uint = 10001
	// 数据库可用时按手机号定位用户，存在则复用其 uid

	if pkg.DB != nil {
		// 声明变量 u model.UserProfile，用于保存运行过程中的数据

		var u model.UserProfile
		// 按手机号查询，命中则采用库中 uid

		if pkg.DB.Where("phone = ?", phone).First(&u).Error == nil {
			// 赋值更新：uid = u.ID

			uid = u.ID
		// 当前代码块结束

		} else {
			// 手机号不存在则创建演示用户，昵称由手机号脱敏生成

			u = model.UserProfile{Phone: phone, Nickname: "用户" + phone[len(phone)-4:], Role: "user", Status: 1}
			// 写库，失败不影响演示登录降级

			if pkg.DB.Create(&u).Error == nil {
				// 赋值更新：uid = u.ID

				uid = u.ID
			// 当前代码块结束

			}
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 调用 IssueToken 签发真实 JWT（含 token_ver，退出登录可即时失效）

	tok, c := IssueToken(uid)
	// 返回 uid, tok, c 给调用方，结束当前函数

	return uid, tok, c
// 当前代码块结束

}
