package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"goApiServer/dbs"
	"goApiServer/lang"
	"goApiServer/structs"
	"goApiServer/tools"
	"strconv"
	"strings"
	"time"
)

// 分页查询, 参数1：第几页，参数2：每页有多少条记录
func SelectUsers(spage, snumber string) (*structs.ReturnUsers, error) {
	// 字符串转int
	page, err01 := strconv.Atoi(spage)
	if err01 != nil {
		beego.Error("字符串转 int 出错：", err01)
		return nil, err01
	}
	// 字符串转int
	number, err02 := strconv.Atoi(snumber)
	if err02 != nil {
		beego.Error("字符串转 int 出错：", err02)
		return nil, err02
	}

	// 小于1都是不合法的, 强制转换为最小的1
	if page < 1 {
		page = 1
	}
	if number < 1 {
		number = 1
	}

	// 定义redis的key, id转string类型
	redisKey := fmt.Sprintf("test:users_number%v_size%v", page, number)

	// 接收数据的变量
	var users structs.ReturnUsers

	// 查询记录
	// 先从缓存查询，没有再从数据库查
	if err := dbs.RedisDB.GetJSON(redisKey, &users); err != nil {
		if strings.Contains(err.Error(), "key not found") {
			// key不存在，就重新查询一次

			// 获取MySQL的查询参数
			limit, offset := tools.DBPage(page, number)

			// sql
			sql := fmt.Sprintf("select * from users limit %v offset %v", limit, offset)

			// 打印查询sql
			beego.Debug("[sql]: ", sql)

			// 接收数据的对象
			var list []structs.User
			err2 := dbs.MysqlDB.Select(&list, sql)
			if err2 != nil {
				beego.Error("Mysql查询出错:", err2)
				return nil, err2
			}

			// 赋值
			users = structs.ReturnUsers{TotalNum: len(list), List: list}

			// 对象存储到缓存
			// 将对象存到缓存中
			if err3 := dbs.RedisDB.SetJSON(redisKey, &users, tools.OneDay); err3 != nil {
				beego.Error("Redis存储出错:", err3)
				return nil, err3
			}

			// 返回查询出来的信息
			return &users, nil
		} else {
			// 打印错误日志
			beego.Error("Redis查询出错:", err)
			// 返回错误信息
			return nil, err
		}
	}

	// 返回缓存中的数据
	return &users, nil
}

// 获取单个用户的信息
func GetUserInfo(id string) (*structs.User, int, error) {
	// 定义redis的key, id转string类型
	redisKey := DataBaseName() + fmt.Sprintf(":user_%v", id)

	// 定义接收数据的对象
	var user structs.User

	// 先从缓存查询，没有再从数据库查
	if err := dbs.RedisDB.GetJSON(redisKey, &user); err != nil {
		if strings.Contains(err.Error(), "key not found") {
			// key不存在，就重新查询一次

			// 预处理SQL语句
			sql := fmt.Sprintf("select * from users where id=%v", id)

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

// 新增用户信息
func AddUser(user *structs.User) (int, error) {

	// 校验用户参数
	if !tools.Validate.CheckUser(user) {
		beego.Error("用户传递的参数中，有不符合要求的参数。详细信息:", tools.ReturnJson(user))
		return lang.RRSErrorInfo.Err1003.Code, lang.RRSErrorInfo.Err1003.ErrorType()
	}

	createAt := time.Now().Unix()

	newPassword, err := tools.PasswordEncryption(user.Name, user.Password, createAt)
	if err != nil {
		// 打印错误日志
		beego.Error("密码加密出错: ", err)
		// 返回错误信息
		return 0, err
	}

	// 预处理SQL语句
	addUserSql := fmt.Sprintf("insert into users (name, password, email, language, role, noticeEnable, noticeLevel, createTime) "+
		"values ('%v', '%v', '%v', '%v', %v, %v, %v, %v)",
		user.Name, newPassword, user.Email, user.Language, user.Role, user.NoticeEnable, user.NoticeLevel, createAt)

	// 打印日志
	beego.Debug("[sql]: " + addUserSql)

	// 执行sql
	_, err = dbs.MysqlDB.Exec(addUserSql)
	if err != nil {
		// 用户名关键字重复错误
		if strings.Contains(err.Error(), "Error 1062") {
			return lang.RRSErrorInfo.Err1001.Code, lang.RRSErrorInfo.Err1001.ErrorType()
		} else {
			// 打印错误日志
			beego.Error("Mysql执行出错: ", err)

			// 返回错误信息
			return 0, err
		}
	}

	// 定义redis的key
	redisKey := DataBaseName() + ":users_"

	// 清空分页数据
	err2 := dbs.CleanRedisPrefix(redisKey)
	if err2 != nil {
		beego.Error(err2)
		return 0, err2
	}

	// 正常情况返回空值
	return 0, nil
}

// 修改用户
func UpdateUser(user *structs.User) (int, error) {
	// 不能缺少用户 ID
	if user.ID == 0 {
		beego.Error("缺少用户ID，无法处理修改用户的请求!")
		return lang.RRSErrorInfo.Err1003.Code, lang.RRSErrorInfo.Err1003.ErrorType()
	}

	// 获取目标用户的信息
	oldUser, code, err := GetUserInfo(strconv.Itoa(user.ID))
	if err != nil {
		// 优先处理，空值报错
		if code == lang.RRSErrorInfo.Err1002.Code {
			return code, err
		} else {
			beego.Error("获取用户数据出错:", err)
			return 0, err
		}
	}

	// 校验更新内容是否合法
	if !tools.Validate.CheckUser(user) {
		beego.Error("用户传递的参数中，有不符合要求的参数。详细信息:", tools.ReturnJson(user))
		return lang.RRSErrorInfo.Err1003.Code, lang.RRSErrorInfo.Err1003.ErrorType()
	}

	// 声明KEY
	redisKey := DataBaseName() + fmt.Sprintf(":user_%v", user.ID)
	// 同一个用户的更新锁和删除锁是相同的，删除的时候不允许更新，更新的时候不能删除
	redisLock := DataBaseName() + fmt.Sprintf(":lock:user_%v", user.ID)

	// 获取redis锁
	ok, err := dbs.RedisDB.Lock(redisLock, time.Minute*10)
	defer dbs.RedisDB.Unlock(redisLock)
	if err != nil {
		beego.Error("Redis加锁失败: ", err)
		return 0, err
	}

	// 加锁成功的，进行进一步操作。
	if ok {
		// 定义一个空对象，用于和传入的值进行比对，得到更新sql
		userNil := structs.User{}

		// 判断用户是否需要更新
		needUpdate := false

		// sql语句
		sql := "update users set"

		// 开始比较
		if user.Name != userNil.Name {
			// 字符串必须有引号，否则sql语句就出错了。
			sql += fmt.Sprintf(" name='%v',", user.Name)
			needUpdate = true
		}
		if user.Password != userNil.Password {
			// 获取到加密后的密码
			newPassword, err := tools.PasswordEncryption(user.Name, user.Password, oldUser.CreateTime)
			if err != nil {
				// 打印错误日志
				beego.Error("密码加密出错: ", err)
				// 返回错误信息
				return 0, err
			}

			// 加密后的密码进行存储
			sql += fmt.Sprintf(" password='%v',", newPassword)
			needUpdate = true
		}
		if user.Email != userNil.Email {
			sql += fmt.Sprintf(" email='%v',", user.Email)
			needUpdate = true
		}
		if user.Language != userNil.Language {
			sql += fmt.Sprintf(" language='%v',", user.Email)
			needUpdate = true
		}
		if user.Role != userNil.Role {
			sql += fmt.Sprintf(" role=%v,", user.Role)
			needUpdate = true
		}
		if user.NoticeEnable != userNil.NoticeEnable {
			sql += fmt.Sprintf(" noticeEnable=%v,", user.NoticeEnable)
			needUpdate = true
		}
		if user.NoticeLevel != userNil.NoticeLevel {
			sql += fmt.Sprintf(" noticeLevel=%v,", user.NoticeLevel)
			needUpdate = true
		}

		// 判断是否需要更新
		if needUpdate {
			// 去掉尾巴上的逗号(只会去掉1个字符，不用但心会删除有用的字符)
			sql = strings.TrimRight(sql, ",")
			// 加上过滤条件
			sql += fmt.Sprintf(" where id=%v", user.ID)
		} else {
			// 所有可更新字段都是空值，那就是什么都不操作。
			return 0, nil
		}

		// 打印日志
		beego.Debug("[sql]: " + sql)

		// 执行sql
		_, err2 := dbs.MysqlDB.Exec(sql)
		if err2 != nil {
			// 打印错误日志
			beego.Error("执行更新 sql 出错:", err2)
			// 返回错误信息
			return 0, err2
		}

		// 清空缓存
		err3 := dbs.RedisDB.Del(redisKey)
		if err3 != nil {
			beego.Error("清空缓存数据失败:", err3)
			return 0, err3
		}

		// 清空分页
		err4 := dbs.CleanRedisPrefix(DataBaseName() + ":users_")
		if err4 != nil {
			beego.Error("清空分页数据失败:", err4)
			return 0, err4
		}

		// 正常情况返回空值
		return 0, nil
	} else {
		// 请求的资源正忙，请稍后再试！
		beego.Error("Redis锁未释放，资源正在被使用，无法处理当前请求。")
		return lang.RRSErrorInfo.Err1004.Code, lang.RRSErrorInfo.Err1004.ErrorType()
	}
}

// 删除用户
func DelUser(id int) (int, error) {
	// 不能缺少用户 ID
	if id == 0 {
		beego.Error("缺少用户ID，无法处理修改用户的请求!")
		return lang.RRSErrorInfo.Err1003.Code, lang.RRSErrorInfo.Err1003.ErrorType()
	}

	// 获取目标用户的信息
	_, code, err := GetUserInfo(strconv.Itoa(id))
	if err != nil {
		// 优先处理，空值报错
		if code == lang.RRSErrorInfo.Err1002.Code {
			return code, err
		} else {
			beego.Error("获取用户数据出错:", err)
			return 0, err
		}
	}

	// 声明KEY
	redisKey := DataBaseName() + fmt.Sprintf(":user_%v", id)
	// 同一个用户的更新锁和删除锁是相同的，删除的时候不允许更新，更新的时候不能删除
	redisLock := DataBaseName() + fmt.Sprintf(":lock:user_%v", id)

	// 获取redis锁
	ok, err := dbs.RedisDB.Lock(redisLock, time.Minute*10)
	defer dbs.RedisDB.Unlock(redisLock)
	if err != nil {
		beego.Error("Redis加锁失败: ", err)
		return 0, err
	}

	// 加锁成功的，进行进一步操作。
	if ok {
		// 预处理SQL语句
		sql := fmt.Sprintf("delete from users where id=%v", id)

		// 打印日志
		beego.Debug("[sql]: ", sql)

		// 执行sql
		_, err2 := dbs.MysqlDB.Exec(sql)
		if err2 != nil {
			// 打印错误日志
			beego.Error("执行更新 sql 出错:", err2)
			// 返回错误信息
			return 0, err2
		}

		// 清空缓存
		err3 := dbs.RedisDB.Del(redisKey)
		if err3 != nil {
			beego.Error("清空缓存数据失败:", err3)
			return 0, err3
		}

		// 清空分页
		err4 := dbs.CleanRedisPrefix(DataBaseName() + ":users_")
		if err4 != nil {
			beego.Error("清空分页数据失败:", err4)
			return 0, err4
		}

		// 正常情况返回空值
		return 0, nil

	} else {
		// 请求的资源正忙，请稍后再试！
		beego.Error("Redis锁未释放，资源正在被使用，无法处理当前请求。")
		return lang.RRSErrorInfo.Err1004.Code, lang.RRSErrorInfo.Err1004.ErrorType()
	}
}
