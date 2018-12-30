package tools

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

//用于获取项目根目录的绝对路径
func GetRootPath() string {
	// 获取项目的上级目录路径
	FileDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return FileDir + "/"
}

// 解析JSON文件: json文件的数据结构 和 接口对象的数据结构 必须要一致！！！
func ParseJsonFile(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	//读取的数据为json格式，需要进行解码
	err2 := json.Unmarshal(data, v)
	if err2 != nil {
		panic(err2)
	}
}
