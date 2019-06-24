package tools

import (
	"fmt"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/scrypt"
)

//MySQL数据库分页计算，传入页数和每页数量，得到sql需要的分页参数
func DBPage(number, size int) (limit int, offset int) {
	// 每页的数量，就是限制
	limit = size
	// 页数-1，再乘以，每页的数量，就等于偏移量
	offset = (number - 1) * size
	// 返回
	return limit, offset
}

// 给密码加密
func PasswordEncryption(username string, password string, createAt int64) (string, error) {
	// 生成盐[用户名、密码、用户创建时间(时间戳类型)]
	salt := StrMD5(username + password + TimestampToLocal(createAt) + "test_salt")

	// 获取加密结果
	dk, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		beego.Error("密码加密出错:", err)
		return "", err
	}

	// 返回加密数据, 因为上面得到的是一个哈希类型的数据，所以要进行格式化处理
	return fmt.Sprintf("%x", dk), nil
}
