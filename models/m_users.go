package models

import (
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"testapi/dbs"
	"testapi/tools"
)

// 用户表结构体,用于接收数据库查询出来的对象，数据类型和数据库尽量保持一致
type User struct {
	Id    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Age   int64  `json:"age" db:"age"`
	Email string `json:"email" db:"email"`
}

// 添加用户时，接收用户传值的对象
type ReceiveUser struct {
	Name  string `json:"name"`
	Age   int64  `json:"age"`
	Email string `json:"email"`
}

// 返回给用户的分页对象
type ReturnUserProducts struct {
	TotalNum int64  `json:"total_num"`
	Users    []User `json:"user_products"`
}

// 查询用户
func SelectUser(id int64) (*User, error) {
	beego.Trace("开始查询用户信息, 查询ID:", id)

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
					if err := dbs.RedisDB.SetJSON(redisKey, nil, tools.OneHour); err != nil {
						beego.Error(err)
						return nil, err
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
			if err := dbs.RedisDB.SetJSON(redisKey, &user, tools.OneDay); err != nil {
				beego.Error(err)
				return nil, err
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
	// DML操作，开启事务
	tx, err := dbs.MysqlDB.Begin()
	if err != nil {
		// 打印错误日志
		beego.Error(err)
		// 返回错误信息
		return err
	}

	// 预处理SQL语句
	addUserSql := "insert into users (name, age, email) values (?, ?, ?)"

	// 打印日志
	beego.Debug("[sql]: "+addUserSql, user.Name, user.Age, user.Email)

	// 执行sql
	_, err2 := tx.Exec(addUserSql, user.Name, user.Age, user.Email)
	if err2 != nil {
		// 打印错误日志
		beego.Error(err2)
		// 回滚
		tx.Rollback()
		// 返回错误信息
		return err2
	}

	//没有问题了，最后一起提交。
	tx.Commit()

	// 正常情况返回空值
	return nil
}

// 修改用户
func UpdateUser(user *User) error {
	// DML操作，开启事务
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
		// 回滚
		tx.Rollback()
		// 返回错误信息
		return err2
	}

	//没有问题了，最后一起提交。
	tx.Commit()

	// 正常情况返回空值
	return nil
}

// 删除用户
func DelUser(id int64) error {
	// DML操作，开启事务
	tx, err := dbs.MysqlDB.Begin()
	if err != nil {
		// 打印错误日志
		beego.Error(err)
		// 返回错误信息
		return err
	}

	// 预处理SQL语句
	updateUserSql := "delete users where id=?"

	// 打印日志
	beego.Debug("[sql]: "+updateUserSql, id)

	// 执行sql
	_, err2 := tx.Exec(updateUserSql, id)
	if err2 != nil {
		// 打印错误日志
		beego.Error(err2)
		// 回滚
		tx.Rollback()
		// 返回错误信息
		return err2
	}

	//没有问题了，最后一起提交。
	tx.Commit()

	// 正常情况返回空值
	return nil
}
