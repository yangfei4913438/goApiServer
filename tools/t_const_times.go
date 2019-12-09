package tools

import (
	"time"
)

// 当前页面，都是时间类型(time.Duration)的值,不做其他用途
const (
	ZeroTime  = time.Second * 0
	OneSecond = time.Second * 1
	OneMinute = OneSecond * 60
	OneHour   = OneMinute * 60
	OneDay    = OneHour * 24
	OneWeek   = OneDay * 7
	OneYear   = OneDay * 365
)

func OneMonth(year int, month int) (int, time.Duration) {
	switch month {
	case 1:
		fallthrough
	case 3:
		fallthrough
	case 5:
		fallthrough
	case 7:
		fallthrough
	case 8:
		fallthrough
	case 10:
		fallthrough
	case 12:
		return 31, OneDay * 31
	case 2:
		if (year % 4) != 0 {
			// 年份不能被4整除，就是平年，二月份为28天
			return 28, OneDay * 28
		} else {
			// 年份可以被4整除，就是闰年，二月份为29天
			return 29, OneDay * 29
		}
	case 4:
		fallthrough
	case 6:
		fallthrough
	case 9:
		fallthrough
	case 11:
		return 30, OneDay * 30
	default:
		return 0, ZeroTime
	}
}

// 获取当前的时区
func GetTimeZone() int {
	// Zone 方法可以获得变量的时区和时区与UTC的偏移秒数
	_, offset := time.Now().Local().Zone()

	// 秒转换为小时返回出去，就是时区
	return offset / 60 / 60
}

// 时区设置
func GetTimeZoneCity(num int) string {
	switch num {
	case -1:
		return "Atlantic/Cape_Verde"
	case -2:
		return "America/Godthab"
	case -3:
		return "America/Bahia"
	case -4:
		return "America/Caracas"
	case -5:
		return "America/Bogota"
	case -6:
		return "America/Belize"
	case -7:
		return "America/Vancouver"
	case -8:
		return "America/Anchorage"
	case -9:
		return "America/Adak"
	case -10:
		return "Pacific/Honolulu"
	case -11:
		return "Pacific/Midway"
	case -12:
		fallthrough
	case 12:
		// 东西十二区是一回事
		return "Pacific/Auckland"
	case 11:
		return "Pacific/Guadalcanal"
	case 10:
		return "Australia/Sydney"
	case 9:
		return "Asia/Tokyo"
	case 8:
		return "Asia/Shanghai"
	case 7:
		return "Asia/Jakarta"
	case 6:
		return "Asia/Dhaka"
	case 5:
		return "Asia/Yekaterinburg"
	case 4:
		return "Asia/Baku"
	case 3:
		return "Europe/Moscow"
	case 2:
		return "Europe/Berlin"
	case 1:
		return "Europe/London"
	case 0:
		return "UTC"
	default:
		// 其他的值，不是时区
		return ""
	}
}
