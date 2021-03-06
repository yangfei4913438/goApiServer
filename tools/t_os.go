package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//用于获取项目根目录的绝对路径
func GetRootPath() (*string, error) {
	// 获取项目的上级目录路径
	FileDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		beego.Error(err)
		return nil, err
	}
	rootPath := FileDir + "/"
	return &rootPath, nil
}

// 解析JSON文件: json文件的数据结构 和 接口对象的数据结构 必须要一致！！！
func ParseJsonFile(filename string, v interface{}) error {
	// ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		beego.Error(err)
		return err
	}

	// 读取的数据为json格式，需要进行解码
	err2 := json.Unmarshal(data, v)
	if err2 != nil {
		beego.Error(err2)
		return err2
	}

	// 无异常就返回nil
	return nil
}

// 执行系统命令【linux】
func ExecBashShell(cmdLine string) (string, error) {
	// 所有变量都要先定义，避免使用直接赋值的变量
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)

	// 生成 cmd
	cmd = exec.Command("/bin/bash", "-c", cmdLine)

	// 执行命令，捕获了子进程的输出(pipe)
	if output, err = cmd.CombinedOutput(); err != nil {
		// 返回错误信息
		return "", err
	} else {
		// 返回输出内容
		return string(output), nil
	}
}

// 设置系统时间【linux】
func SetLinuxSystemTime(tt int64) error {
	// 时间戳转成本地时间
	localTime := TimestampToLocal(tt)

	// 组合命令语句
	cmdLine := fmt.Sprintf("date -s '%v'", localTime)

	// 执行命令
	if _, err := ExecBashShell(cmdLine); err != nil {
		return err
	} else {
		// 写入CMOS
		return syncToCOMS()
	}
}

// 设置系统时区【linux】
func SetLinuxTimeZone(tz int) error {
	// 先获取时区对应的城市
	city := GetTimeZoneCity(tz)

	// 判断值是否正确
	if len(city) == 0 {
		beego.Error("时区错误，取值范围应该是-12到12之间的整数。但是，接收到的参数是:", tz)
		// 1003 表示参数错误
		return errors.New("1003")
	}

	// 组合命令语句
	cmdLine := fmt.Sprintf("timedatectl set-timezone '%v'", city)

	// 执行命令
	if _, err := ExecBashShell(cmdLine); err != nil {
		return err
	} else {
		// 写入CMOS
		return syncToCOMS()
	}
}

// 将设置写入硬件时钟【linux】
func syncToCOMS() error {
	// 执行命令
	if _, err := ExecBashShell("hwclock -w"); err != nil {
		return err
	} else {
		return nil
	}
}

// 使用shell方式从系统中获取时区【MAC,linux限定】
func GetShellTimezone() (int, error) {
	res, err := ExecBashShell("date -R")
	if err != nil {
		fmt.Println("【异常】服务器执行shell出错:", err)
		return 500, err
	}
	// 去掉换行符
	str := strings.Replace(res, "\n", "", 1)

	// 分割字符串
	arr := strings.Split(str, " ")

	// 获取索引
	index := len(arr) - 1

	// 取到时区的字符串
	str = arr[index]

	// 判断并返回
	switch str {
	case "+0000":
		return 0, nil
	case "-0100":
		return -1, nil
	case "-0200":
		return -2, nil
	case "-0300":
		return -3, nil
	case "-0400":
		return -4, nil
	case "-0500":
		return -5, nil
	case "-0600":
		return -6, nil
	case "-0700":
		return -7, nil
	case "-0800":
		return -8, nil
	case "-0900":
		return -9, nil
	case "-1000":
		return -10, nil
	case "-1100":
		return -11, nil
	case "-1200":
		fallthrough
	case "+1200":
		return 12, nil
	case "+1100":
		return 11, nil
	case "+1000":
		return 10, nil
	case "+0900":
		return 9, nil
	case "+0800":
		return 8, nil
	case "+0700":
		return 7, nil
	case "+0600":
		return 6, nil
	case "+0500":
		return 5, nil
	case "+0400":
		return 4, nil
	case "+0300":
		return 3, nil
	case "+0200":
		return 2, nil
	case "+0100":
		return 1, nil
	default:
		return 500, errors.New("undefined timezone! ")
	}
}
