package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"testapi/dbs"
	"testapi/lang"
	"testapi/tools"
	"time"
)

// 用户表结构体,用于接收数据库查询出来的对象，数据类型和数据库尽量保持一致
type User struct {
	Id          int64 `json:"id" db:"id"`
	ReceiveUser       // 直接放struct对象，就等于直接把这个struct的属性放到这里，不是新增一个对象，而是新增三个属性（针对当前操作）
}

// 添加用户时，接收用户传值的对象
type ReceiveUser struct {
	Name  string `json:"name" db:"name"`
	Age   int64  `json:"age" db:"age"`
	Email string `json:"email" db:"email"`
}

// 返回给用户的分页对象
type ReturnUsers struct {
	TotalNum int    `json:"total_num"`
	List     []User `json:"list"`
}

// 分页查询, 参数1：第几页，参数2：每页有多少条记录
func SelectUsers(spage, snumber string) (*ReturnUsers, error) {
	// 字符串转int
	page, err01 := strconv.Atoi(spage)
	if err01 != nil {
		beego.Error(err01)
		return nil, err01
	}
	// 字符串转int
	number, err02 := strconv.Atoi(snumber)
	if err02 != nil {
		beego.Error(err02)
		return nil, err02
	}

	// 定义redis的key, id转string类型
	redisKey := fmt.Sprintf("test:users_number%v_size%v", page, number)

	// 小于1都是不合法的, 强制转换为最小的1
	if page < 1 {
		page = 1
	}
	if number < 1 {
		number = 1
	}

	// 接收数据的变量
	var users *ReturnUsers

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
			beego.Debug(sql)
			// 接收数据的对象
			var list []User
			err2 := dbs.MysqlDB.Select(&list, sql)
			if err2 != nil {
				beego.Error(err2)
				return nil, err2
			}

			// 赋值
			users = &ReturnUsers{len(list), list}

			// 对象存储到缓存
			// 将对象存到缓存中
			if err3 := dbs.RedisDB.SetJSON(redisKey, users, tools.OneDay); err3 != nil {
				beego.Error(err3)
				return nil, err3
			}

			// 返回查询出来的信息
			return users, nil
		} else {
			// 打印错误日志
			beego.Error(err)
			// 返回错误信息
			return nil, err
		}
	}

	// 返回缓存中的数据
	return users, nil
}

// 查询用户
func SelectUser(id int64) (resObj *User, resErr error) {
	// 多语言打印：现在开始查询用户的信息。查询ID:
	beego.Trace(lang.CurrLang.Models.Users.SelectInfo01, id)
	// 多语言打印：查询用户信息操作已完成。
	defer beego.Trace(lang.CurrLang.Models.Users.SelectInfo02)

	// 定义redis的key, id转string类型
	redisKey := fmt.Sprintf("test:user_%v", id)

	// 定义接收数据的对象
	var user User

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
					beego.Trace("查询结果为空值!")

					// 将空值添加到缓存, 有效期1小时
					if err1 := dbs.RedisDB.SetJSON(redisKey, nil, tools.OneHour); err1 != nil {
						beego.Error(err1)
						return nil, err1
					}
					// 返回空值
					return nil, nil
				} else {
					// 打印错误日志
					beego.Error(err)
					// 返回错误信息
					return nil, err
				}
			}

			// 将结果添加到缓存
			if err2 := dbs.RedisDB.SetJSON(redisKey, &user, tools.OneDay); err2 != nil {
				// 打印错误日志
				beego.Error(err2)
				// 返回错误信息
				return nil, err2
			}

			// 返回结果给用户
			return &user, nil
		}
	}

	// 空值返回用户信息
	if user.Id == 0 {
		// 因为正常情况下，ID是从1开始的。0就表示读取出来的值是空值
		return nil, nil
	} else {
		return &user, nil
	}
}

// 新增用户信息
func AddUser(user *ReceiveUser) error {
	// 这里的sql写法不一样，主要是为了展示mysql模块api传参的用法

	// 预处理SQL语句
	addUserSql := "insert into users (name, age, email) values (?, ?, ?)"

	// 打印日志
	beego.Debug("[sql]: "+addUserSql, user.Name, user.Age, user.Email)

	// 执行sql
	_, err := dbs.MysqlDB.Exec(addUserSql, user.Name, user.Age, user.Email)
	if err != nil {
		// 打印错误日志
		beego.Error(err)
		// 返回错误信息
		return err
	}

	// 清空分页
	err2 := dbs.CleanRedisPrefix("test:users_")
	if err2 != nil {
		beego.Error(err2)
		return err2
	}

	// 正常情况返回空值
	return nil
}

// 修改用户
func UpdateUser(user *User) error {
	// 声明KEY
	redisKey := fmt.Sprintf("test:user_%v", user.Id)
	// 同一个用户的更新锁和删除锁是相同的，删除的时候不允许更新，更新的时候不能删除
	redisLock := fmt.Sprintf("test:lock:user_%v", user.Id)

	// 获取redis锁
	ok, err := dbs.RedisDB.Lock(redisLock, time.Minute*10)
	defer dbs.RedisDB.Unlock(redisLock)
	if err != nil {
		beego.Error(err)
		return err
	}

	// 加锁成功的，进行进一步操作。
	if ok {
		// 开启事务, 事务功能的demo
		// 正常情况下，如果只有单条SQL语句执行，请不要使用事务！这里只是展示事务的用法，所以才会出现事务。
		tx, err := dbs.MysqlDB.Begin()
		if err != nil {
			// 打印错误日志
			beego.Error(err)
			// 返回错误信息
			return err
		}

		// 定义一个空对象，用于和传入的值进行比对，得到更新sql
		userNil := User{}

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
		if user.Age != userNil.Age {
			// 数字类型，不需要加上引号
			sql += fmt.Sprintf(" age=%v,", user.Age)
			needUpdate = true
		}
		if user.Email != userNil.Email {
			sql += fmt.Sprintf(" email='%v',", user.Email)
			needUpdate = true
		}

		// 判断是否需要更新
		if needUpdate {
			// 去掉尾巴上的逗号(只会去掉1个字符，不用但心会删除有用的字符)
			sql = strings.TrimRight(sql, ",")
			// 加上过滤条件
			sql += fmt.Sprintf(" where id=%v", user.Id)
		} else {
			// 所有可更新字段都是空值，那就是什么都不操作。
			return nil
		}

		// 打印日志
		beego.Debug("[sql]: " + sql)

		// 执行sql
		_, err2 := tx.Exec(sql)
		if err2 != nil {
			// 打印错误日志
			beego.Error(err2)
			// 回滚, 这里失败了，也是要打印的
			rollErr := tx.Rollback()
			if rollErr != nil {
				// 这里的错误不需要返回了，因为不属于主要错误，返回主要错误，让查询的时候，自然就会查到这里了。
				beego.Error(rollErr.Error())
			}
			// 返回错误信息
			return err2
		}

		//没有问题了，最后一起提交。
		cmtErr := tx.Commit()
		if cmtErr != nil {
			beego.Error(cmtErr.Error())
			// 提交出错肯定是要返回错误信息的。
			return cmtErr
		}

		// 清空缓存
		err3 := dbs.RedisDB.Del(redisKey)
		if err3 != nil {
			beego.Error(err3)
			return err3
		}

		// 清空分页
		err4 := dbs.CleanRedisPrefix("test:users_")
		if err4 != nil {
			beego.Error(err4)
			return err4
		}

		// 正常情况返回空值
		return nil
	} else {
		return errors.New("资源被锁定，无法更新数据，请稍后再试! ")
	}
}

// 删除用户
func DelUser(id int64) error {
	// 声明KEY
	redisKey := fmt.Sprintf("test:user_%v", id)
	// 同一个用户的更新锁和删除锁是相同的，删除的时候不允许更新，更新的时候不能删除
	redisLock := fmt.Sprintf("test:lock:user_%v", id)

	// 获取redis锁
	ok, err := dbs.RedisDB.Lock(redisLock, time.Minute*10)
	defer dbs.RedisDB.Unlock(redisLock)
	if err != nil {
		beego.Error(err)
		return err
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
			beego.Error(err2)
			// 返回错误信息
			return err2
		}

		// 清空缓存
		err3 := dbs.RedisDB.Del(redisKey)
		if err3 != nil {
			beego.Error(err3)
			return err3
		}

		// 清空分页
		err4 := dbs.CleanRedisPrefix("test:users_")
		if err4 != nil {
			beego.Error(err4)
			return err4
		}

		// 正常情况返回空值
		return nil
	} else {
		return errors.New("资源被锁定，无法删除，请稍后再试! ")
	}
}
