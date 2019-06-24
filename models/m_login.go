package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/mojocn/base64Captcha"
	"goApiServer/dbs"
	"goApiServer/lang"
	"goApiServer/structs"
	"goApiServer/tools"
	"strings"
)

// 数据库名称
func DataBaseName() string {
	return beego.AppConfig.String("mysql_db")
}

// 创建token
func getToken() string {

	res := ""

	// 取 64 个随机字符
	for i := 0; i < 64; i++ {
		res += tools.RandStr()
	}

	// 返回 token
	return tools.StrSHA256(res)
}

// 获取图形验证码：数字
func GetImageCode() (captchaId string, imageCode string) {
	// 数字验证码配置
	var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 4,
	}

	// 创建数字验证码.
	// GenerateCaptcha 第一个参数为空字符串,包会自动在服务器一个随机种子给你产生随机uuid.
	captchaId, capD := base64Captcha.GenerateCaptcha("", configD)

	// 以base64编码
	imageCode = base64Captcha.CaptchaWriteToBase64Encoding(capD)

	beego.Notice("获取到的验证码信息:", captchaId, imageCode)
	return captchaId, imageCode
}

// 根据用户名获取用户信息【因为用户名是唯一的，所以可以查询】
func getUserInfoByName(name string) (*structs.User, int, error) {

	// 定义redis的key, id转string类型
	redisKey := DataBaseName() + fmt.Sprintf(":user_%v", name)

	// 定义接收数据的对象
	var user structs.User

	// 先从缓存查询，没有再从数据库查
	if err := dbs.RedisDB.GetJSON(redisKey, &user); err != nil {
		if strings.Contains(err.Error(), "key not found") {
			// key不存在，就重新查询一次

			// 预处理SQL语句
			sql := fmt.Sprintf("select * from users where name='%v'", name)

			// 打印日志
			beego.Debug("[sql]:", sql)
			err := dbs.MysqlDB.Get(&user, sql)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					// 将空对象添加到缓存, 有效期1小时
					if err1 := dbs.RedisDB.SetJSON(redisKey, &user, tools.OneHour); err1 != nil {
						beego.Error("Redis存储出错:", err1)
						// 存储出错，返回错误信息
						return nil, 0, err
					}

					// 空对象返回错误代码
					return nil, lang.RRSErrorInfo.Err1002.Code, lang.RRSErrorInfo.Err1002.ErrorType()

				} else {
					// 打印错误日志
					beego.Error("Mysql查询出错:", err)
					// 返回错误信息
					return nil, 0, err
				}
			}

			// 将结果添加到缓存
			if err2 := dbs.RedisDB.SetJSON(redisKey, &user, tools.OneDay); err2 != nil {
				// 打印错误日志
				beego.Error("Redis存储出错:", err2)
				// 返回错误信息
				return nil, 0, err2
			}

			// 返回结果给用户
			return &user, 0, nil
		} else {
			// 打印错误日志
			beego.Error("Redis查询出错:", err)
			// 返回错误信息
			return nil, 0, err
		}
	}

	// 空值返回用户信息
	if user.ID == 0 {
		// 空对象返回错误代码
		return nil, lang.RRSErrorInfo.Err1002.Code, lang.RRSErrorInfo.Err1002.ErrorType()
	} else {
		return &user, 0, nil
	}
}

// 用户登出
func Logout(token string) error {
	// 生成 token key
	tokenKey := "user_token:" + token

	// 执行删除操作
	err := dbs.RedisDB.Del(tokenKey)
	if err != nil {
		beego.Error("Redis 删除用户 token，操作出错:", err)
		return err
	}

	return nil
}

// 用户登录
func Login(user *structs.User) *structs.LoginReturn {
	var ServerErrorCode = 0

	// 不能缺少登录账户或登录密码
	if user.Name == "" || user.Password == "" {
		beego.Error("无效的登录账户或登录密码，无法处理用户的登陆请求!")
		return &structs.LoginReturn{
			ErrCode: lang.RRSErrorInfo.Err1003.Code,
			ErrMsg:  lang.RRSErrorInfo.Err1003.Message,
		}
	}
	// 用户的验证码查询
	userCheckKey := "user_captcha:" + user.Name
	ok, err := dbs.RedisDB.Exists(userCheckKey)
	if err != nil {
		beego.Error("Redis获取用户是否有验证码时，查询出错了:", err)
		return &structs.LoginReturn{
			ErrCode: ServerErrorCode,
			ErrMsg:  err.Error(),
		}
	}
	if ok {
		if user.CaptchaId == "" || user.ImgCode == "" {
			beego.Error("用户缺少验证码，请提交验证码！")
			return &structs.LoginReturn{
				ErrCode: lang.RRSErrorInfo.Err1003.Code,
				ErrMsg:  lang.RRSErrorInfo.Err1003.Message,
			}
		}
		var checkID string
		if err := dbs.RedisDB.GetJSON(userCheckKey, &checkID); err != nil {
			if !strings.Contains(err.Error(), "key not found") {
				beego.Error("Redis获取用户校验 ID 出错:", err)
				return &structs.LoginReturn{
					ErrCode: ServerErrorCode,
					ErrMsg:  err.Error(),
				}
			}
		}
		if checkID != user.CaptchaId {
			beego.Error("用户提交的验证码校验 ID 不存在，初步判断属于恶意测试！")
			return &structs.LoginReturn{
				ErrCode: lang.RRSErrorInfo.Err1003.Code,
				ErrMsg:  lang.RRSErrorInfo.Err1003.Message,
			}
		} else {
			// 比对验证码
			if !base64Captcha.VerifyCaptcha(user.CaptchaId, user.ImgCode) {
				beego.Error("验证码不匹配!")
				return &structs.LoginReturn{
					ErrCode: lang.RRSErrorInfo.Err1008.Code,
					ErrMsg:  lang.RRSErrorInfo.Err1008.Message,
				}
			}
		}
	}

	// 登录拒绝 key
	loginDisableKey := "login_disable:" + user.Name
	ok, err = dbs.RedisDB.Exists(loginDisableKey)
	if err != nil {
		beego.Error("Redis获取用户是否被禁止登录时，查询出错了:", err)
		return &structs.LoginReturn{
			ErrCode: ServerErrorCode,
			ErrMsg:  err.Error(),
		}
	}
	if ok {
		beego.Error("该用户错误登录错误次数过多，暂时不允许登录！")
		return &structs.LoginReturn{
			ErrCode: lang.RRSErrorInfo.Err1007.Code,
			ErrMsg:  lang.RRSErrorInfo.Err1007.Message,
		}
	}

	// 登录错误次数 KEY
	loginErrorKey := "login_error_count:" + user.Name
	var errCnt int
	if err := dbs.RedisDB.GetJSON(loginErrorKey, &errCnt); err != nil {
		if !strings.Contains(err.Error(), "key not found") {
			beego.Error("Redis获取用户登录错误次数出错:", err)
			return &structs.LoginReturn{
				ErrCode: ServerErrorCode,
				ErrMsg:  err.Error(),
			}
		}
	}
	// 获取登录错误次数，最大允许出错 10 次。登录错误达到10次，5 分钟内不允许再次登录！
	if errCnt >= 10 {
		// 这里的 value, 随便写个值就行了。我这里是对 key 的说明。检测的时候，只检测是否有key。
		if err = dbs.RedisDB.SetJSON(loginDisableKey, "user disable login!", tools.OneMinute*5); err != nil {
			// 打印错误日志
			beego.Error("Redis存储出错:", err)
			return &structs.LoginReturn{
				ErrCode: ServerErrorCode,
				ErrMsg:  err.Error(),
			}
		}

		// 清除记录的错误信息，让用户 5 分钟内不能登录，错误次数就没用了。
		// 登录错误次数
		_ = dbs.RedisDB.Del(loginErrorKey)
		// 用户的验证码查询key
		_ = dbs.RedisDB.Del(userCheckKey)

		// 身份验证失败次数过多，请 5 分钟后再尝试登录。
		beego.Error("登录错误达到10次，5 分钟内不允许再次登录！")
		return &structs.LoginReturn{
			ErrCode: lang.RRSErrorInfo.Err1006.Code,
			ErrMsg:  lang.RRSErrorInfo.Err1006.Message,
		}
	}

	// 获取目标用户的信息
	oldUser, code, err := getUserInfoByName(user.Name)
	if err != nil {
		// 优先处理，空值报错
		if code == lang.RRSErrorInfo.Err1002.Code {
			return &structs.LoginReturn{
				ErrCode: code,
				ErrMsg:  err.Error(),
			}
		} else {
			beego.Error("获取用户数据出错:", err)
			return &structs.LoginReturn{
				ErrCode: ServerErrorCode,
				ErrMsg:  err.Error(),
			}
		}
	}

	// 获取到加密后的密码
	newPassword, err := tools.PasswordEncryption(user.Name, user.Password, oldUser.CreateTime)
	if err != nil {
		// 打印错误日志
		beego.Error("密码加密出错: ", err)
		// 返回错误信息
		return &structs.LoginReturn{
			ErrCode: ServerErrorCode,
			ErrMsg:  err.Error(),
		}
	}

	// 比对密码是否一致
	if oldUser.Password == newPassword {
		tk := getToken()
		// token 存储的 key 需要加上前缀，用于 token 存在校验
		tokenKey := "user_token:" + tk
		// 因为用户资料是可变的，用户 ID 是不变的，所以这里直接存用户 ID 是比较稳妥的。
		if err = dbs.RedisDB.SetJSON(tokenKey, oldUser.ID, tools.OneHour*2); err != nil {
			// 打印错误日志
			beego.Error("Redis存储出错:", err)
			return &structs.LoginReturn{
				ErrCode: ServerErrorCode,
				ErrMsg:  err.Error(),
			}
		} else {
			// 清除记录的错误信息， 只要成功一次，就清除失败信息。
			// 登录错误次数
			_ = dbs.RedisDB.Del(loginErrorKey)
			// 用户的验证码查询key
			_ = dbs.RedisDB.Del(userCheckKey)

			// 返回 token 给用户
			return &structs.LoginReturn{
				Token:    tk,
				UserInfo: oldUser,
			}
		}
	} else {
		if errCnt >= 3 { // 失败达到三次，登录需要验证码。
			captchaId, imageCode := GetImageCode()

			// 更新错误次数
			errCnt += 1
			if err := dbs.RedisDB.SetJSON(loginErrorKey, errCnt, tools.OneHour*1); err != nil {
				// 打印错误日志
				beego.Error("Redis存储出错:", err)
				return &structs.LoginReturn{
					ErrCnt: ServerErrorCode,
					ErrMsg: err.Error(),
				}
			}
			// 记录当前用户已经存在验证码
			userCheckKey := "user_captcha:" + user.Name
			if err := dbs.RedisDB.SetJSON(userCheckKey, captchaId, tools.OneHour*1); err != nil {
				// 打印错误日志
				beego.Error("Redis存储出错:", err)
				return &structs.LoginReturn{
					ErrCnt: ServerErrorCode,
					ErrMsg: err.Error(),
				}
			}
			// 返回错误信息
			return &structs.LoginReturn{
				CaptchaId: captchaId,
				ImgCode:   imageCode,
				ErrCnt:    errCnt,
				ErrCode:   lang.RRSErrorInfo.Err1005.Code,
				ErrMsg:    lang.RRSErrorInfo.Err1005.Message,
			}
		} else {
			// 更新错误次数
			errCnt += 1
			if err := dbs.RedisDB.SetJSON(loginErrorKey, errCnt, tools.OneHour*1); err != nil {
				// 打印错误日志
				beego.Error("Redis存储出错:", err)
				return &structs.LoginReturn{
					ErrCnt: ServerErrorCode,
					ErrMsg: err.Error(),
				}
			}
			// 返回错误信息
			return &structs.LoginReturn{
				ErrCnt:  errCnt,
				ErrCode: lang.RRSErrorInfo.Err1005.Code,
				ErrMsg:  lang.RRSErrorInfo.Err1005.Message,
			}
		}
	}
}
