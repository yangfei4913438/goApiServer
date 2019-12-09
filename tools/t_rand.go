package tools

import (
	"math/rand"
	"time"
)

//生成指定范围内的随机数【0到x，不包含x】
func RandInt(x int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(x)
}

//生成随机字符
func RandStr() string {
	src := [...]string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c",
		"d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p",
		"q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C",
		"D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P",
		"Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "~", "!", "@",
		"#", "$", "%", "^", "&", "*", "(", ")", "-", "_", "=", "+", "[",
		"]", "{", "}", "|", "<", ">", "?", "/", ".", ",", ";", ":"}

	return src[RandInt(len(src))]
}

// 生成指定范围内的随机数
func RandRangeInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}
