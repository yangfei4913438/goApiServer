package tools

import (
	"errors"
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

func OneMonth(year int, month int) (int, time.Duration, error) {
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
		return 31, OneDay * 31, nil
	case 2:
		if (year % 4) != 0 {
			// 年份不能被4整除，就是平年，二月份为28天
			return 28, OneDay * 28, nil
		} else {
			// 年份可以被4整除，就是闰年，二月份为29天
			return 29, OneDay * 29, nil
		}
	case 4:
		fallthrough
	case 6:
		fallthrough
	case 9:
		fallthrough
	case 11:
		return 30, OneDay * 30, nil
	default:
		return 0, ZeroTime, errors.New("不是正常的月份值")
	}
}
