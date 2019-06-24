package tools

import (
	"github.com/astaxie/beego"
	"runtime"
	"time"
)

// 两个时间的时间差，返回结果为秒
func TimeDifference(oldTime, newTime string) int {
	//传入的时间，必须是规定的格式 "2006-01-02 15:04:05"
	//用于计算2个时间之间有多少秒的差别

	old_time, _ := time.Parse("2006-01-02 15:04:05", oldTime)
	new_time, _ := time.Parse("2006-01-02 15:04:05", newTime)

	if old_time.Unix() >= new_time.Unix() {
		subTimes := old_time.Sub(new_time).Seconds()
		return int(subTimes)
	} else {
		subTimes := new_time.Sub(old_time).Seconds()
		return int(subTimes)
	}
}

// 返回当前时间
func CurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 返回当前时间的时间戳
func CurrentTimestamp() int64 {
	return time.Now().Unix()
}

// 根据时间戳返回本地时间
func TimestampToLocal(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

// 定时器函数
func timeTicker(f func(), t time.Duration) {
	// 将线程绑定到系统线程
	runtime.LockOSThread()
	ticker := time.NewTicker(t)
	defer func() {
		ticker.Stop()
		runtime.UnlockOSThread()
		// panic处理，打印错误信息
		if err := recover(); err != nil {
			beego.Error("定时器故障! 捕获到了panic错误：", err)
		}
	}()

	// 遍历
	for _ = range ticker.C {
		f()
	}
}
