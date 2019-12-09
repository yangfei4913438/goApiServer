package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/scrypt"
)

func StrMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 给密码加密，这里推荐使用用户ID。后面有空的时候我会更新一下。
// 这种方式加密的弊端是用户名不能修改了，否则就无法登录了。正常的用法是使用用户ID
func PasswordEncryption(username string, password string, createAt int64) (string, error) {
	// 生成盐[用户名、密码、用户创建时间(时间戳类型)]
	salt := StrMD5(username + password + fmt.Sprintf("%v", createAt) + "test_salt")

	// 获取加密结果
	dk, err := scrypt.Key([]byte(password), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		beego.Error("密码加密出错:", err)
		return "", err
	}

	// 返回加密数据, 因为上面得到的是一个哈希类型的数据，所以要进行格式化处理
	return fmt.Sprintf("%x", dk), nil
}

func main() {
	pass, _ := PasswordEncryption("admin", "123456", 1573639897)
	fmt.Println(pass)
}

// 直接利用这种方式，可以直接计算出初始密码。 然后作为默认密码插入到数据库里面去。
// 当然，我这里只是一个抛砖引玉，你完全可以用自己的方法去实现。
