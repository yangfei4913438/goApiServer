package models

import (
	"github.com/mojocn/base64Captcha"
	"goApiServer/structs"
)

// 返回图片验证码
func GetImgCodeUrl() *structs.LoginReturn {

	//4位长度的验证码
	captchaId, imageCode := GetImageCode()

	// 返回错误信息
	return &structs.LoginReturn{
		CaptchaId: captchaId,
		ImgCode:   imageCode,
	}
}

// 验证码校验码
func CheckImgCode(id string, code string) bool {
	// 返回校验结果
	return base64Captcha.VerifyCaptcha(id, code)
}
