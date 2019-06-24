package tools

import (
	"github.com/astaxie/beego"
	"goApiServer/structs"
	"regexp"
	"strings"
)

type validator struct{}

var Validate validator

// 验证用户名的长度
func (v *validator) CheckUserName(name string) bool {
	// 用户名的长度必须大于等于 2个字符，且小于等于 20个字符。
	res := len(name) >= 2 && len(name) <= 20
	if res {
		return true
	} else {
		beego.Error("用户名的长度必须大于等于 2个字符，且小于等于 20个字符！")
		return false
	}
}

// 验证密码
func (v *validator) CheckPassword(pwd string) bool {
	// 规则 1：必须含有数字 0-9
	reg1 := regexp.MustCompile(`[0-9]+`)

	// 规则 2：必须含有大写字母 A-Z
	reg2 := regexp.MustCompile(`[A-Z]+`)

	// 规则 3：必须含有小写字母 a-z
	reg3 := regexp.MustCompile(`[a-z]+`)

	// 规则 4：必须含有指定的特殊字符
	reg4 := regexp.MustCompile(`[-=[;,./~!@#$%^*()_+}{:?]+`)

	// 规则 5：密码长度必须在 6-20 位之间
	reg5 := regexp.MustCompile(`^[\s\S]{6,20}$`)

	// 必须同时满足 5 种规则
	res := reg1.MatchString(pwd) && reg2.MatchString(pwd) && reg3.MatchString(pwd) && reg4.MatchString(pwd) && reg5.MatchString(pwd)

	if res {
		return true
	} else {
		beego.Error("密码必须含有至少: 一个数字、一个小写字母、一个大写字母、一个特殊符号，并且长度在 6-20 位之间！")
		return false
	}
}

// 验证电子邮件
func (v *validator) CheckEmail(email string) bool {
	if email == "" {
		// 没有填写 email 的情况下，判断通过，因为允许为空
		return true
	}
	reg := regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
	res := reg.MatchString(email)
	if res {
		return true
	} else {
		beego.Error("电子邮件的格式错误!")
		return false
	}
}

// 验证语言
func (v *validator) CheckLanguage(lang string) bool {
	// 先将语言转换为小写字母
	switch strings.ToLower(lang) {
	case "zh-cn":
		return true
	case "zh-tw":
		return true
	case "en":
		return true
	default:
		beego.Error("未知的语言类型!")
		return false
	}
}

// 验证用户角色
func (v *validator) CheckRole(role int) bool {
	switch role {
	case 1:
		// 普通用户
		return true
	case 9:
		// 管理员
		return true
	default:
		beego.Error("未知的用户角色类型!")
		return false
	}
}

// 验证邮件通知开关
func (v *validator) CheckNotification(enable int) bool {
	switch enable {
	case 0:
		// 禁用
		return true
	case 1:
		// 启用
		return true
	default:
		beego.Error("未知的开关标识!")
		return false
	}
}

// 验证邮件通知等级
func (v *validator) CheckNotificationLevel(level int) bool {
	switch level {
	case -1:
		// warning
		return true
	case -2:
		// error
		return true
	case -3:
		// fatal
		return true
	default:
		beego.Error("未知的邮件通知等级!")
		return false
	}
}

// 直接验证用户
func (v *validator) CheckUser(user *structs.User) bool {
	return v.CheckUserName(user.Name) &&
		v.CheckPassword(user.Password) &&
		v.CheckEmail(user.Email) &&
		v.CheckLanguage(user.Language) &&
		v.CheckRole(user.Role) &&
		v.CheckNotification(user.NoticeEnable) &&
		v.CheckNotificationLevel(user.NoticeLevel)
}
