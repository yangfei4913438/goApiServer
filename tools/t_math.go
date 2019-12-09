package tools

import (
	"strconv"
)

//单位转换: 传入值的单位(B)
func MathChangeUnit(unit int64) string {
	m := int64(1024)
	switch {
	case unit < m:
		return Int64ToStr(unit) + "B"
	case m <= unit && unit < m*m:
		tmp := float64(unit) / 1024
		return Float64ToStr(tmp) + "KB"
	case m*m <= unit && unit < m*m*m:
		tmp := float64(unit) / 1024 / 1024
		return Float64ToStr(tmp) + "MB"
	case m*m*m <= unit && unit < m*m*m*m:
		tmp := float64(unit) / 1024 / 1024 / 1024
		return Float64ToStr(tmp) + "GB"
	default:
		tmp := float64(unit) / 1024 / 1024 / 1024 / 1024
		return Float64ToStr(tmp) + "TB"
	}
}

func Int64ToStr(x int64) string {
	return strconv.FormatInt(x, 10)
}
func Float64ToStr(x float64) string {
	return strconv.FormatFloat(x, 'f', 2, 64)
}
func IntToStr(x int) string {
	return strconv.Itoa(x)
}
func StrToInt(x string) (int, error) {
	return strconv.Atoi(x)
}
