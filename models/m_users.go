package models

import (
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"testapi/dbs"
	"testapi/lang"
	"testapi/tools"
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
type ReturnUserProducts struct {
	TotalNum int64  `json:"total_num"`
	Users    []User `json:"user_products"`
}

// 查询用户
func SelectUser(id int64) (resObj *User, resErr error) {
	// 多语言打印：现在开始查询用户的信息。查询ID:
	beego.Trace(lang.CurrLang.Models.Users.SelectInfo01, id)
	// 多语言打印：查询用户信息操作已完成。
	defer beego.Trace(lang.CurrLang.Models.Users.SelectInfo02)

	// 定义redis的key, id转string类型
	redisKey := "test:user_" + strconv.FormatInt(id, 10)

	// 定义接收数据的对象
	var user User

	// 先从缓存查询，没有再从数据库查
	if err := dbs.RedisDB.GetJSON(redisKey, &user); err != nil {
		if strings.Contains(err.Error(), "key not found") {
			// key不存在，就重新查询一次

			// 预处理SQL语句
			selectSql := "select * from users where id=? limit 1"

			// 打印日志
			beego.Debug("[sql]: "+selectSql, id)
			err := dbs.MysqlDB.Get(&user, selectSql, id)
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

	// 正常情况返回空值
	return nil
}

// 修改用户
func UpdateUser(user *User) error {
	// 开启事务, 事务功能的demo
	// 正常情况下，如果只有单条SQL语句执行，请不要使用事务！这里只是展示事务的用法，所以才会出现事务。
	tx, err := dbs.MysqlDB.Begin()
	if err != nil {
		// 打印错误日志
		beego.Error(err)
		// 返回错误信息
		return err
	}

	// 预处理SQL语句
	updateUserSql := "update users set name=?, age=?, email=? where id=?"

	// 打印日志
	beego.Debug("[sql]: "+updateUserSql, user.Name, user.Age, user.Email, user.Id)

	// 执行sql
	_, err2 := tx.Exec(updateUserSql, user.Name, user.Age, user.Email, user.Id)
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

	// 正常情况返回空值
	return nil
}

// 删除用户
func DelUser(id int64) error {
	// 预处理SQL语句
	updateUserSql := "delete users where id=?"

	// 打印日志
	beego.Debug("[sql]: "+updateUserSql, id)

	// 执行sql
	_, err := dbs.MysqlDB.Exec(updateUserSql, id)
	if err != nil {
		// 打印错误日志
		beego.Error(err)
		// 返回错误信息
		return err
	}

	// 正常情况返回空值
	return nil
}
