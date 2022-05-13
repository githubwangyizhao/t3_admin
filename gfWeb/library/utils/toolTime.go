package utils

import (
	"fmt"
	"time"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

//获取今日0点时间戳
func GetTodayZeroTimestamp() int {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return int(tm1.Unix())
}

// 时间格式化成文件名
func TimeFormFileName() string {
	currTime := gtime.Timestamp()
	return TimeIntForm(currTime, "20060102-150405")
}

// 时间戳格式化
func TimeInt64FormDefault(timeInt int64) string {
	return TimeIntForm(timeInt, "2006-01-02 15:04:05")
}
func TimeIntFormDefault(timeInt int) string {
	return TimeInt64FormDefault(gconv.Int64(timeInt))
}
func TimeIntForm(timeInt int64, formStr string) string {
	tm := time.Unix(timeInt, 0)
	return tm.Format(formStr)
}

func FormatTime(sec int) string {
	h := sec / (60 * 60)
	m := sec % (60 * 60) / 60
	return fmt.Sprintf("%d:%d", h, m)
}

func FormatTime64Second(sec int64) string {
	return FormatTimeSecond(gconv.Int(sec))
}
func FormatTimeSecond(sec int) string {
	str := ""
	d := sec / (86400)
	h := (sec % 86400) / (3600)
	m := (sec % 3600) / 60
	s := sec % 60
	if d > 0 {
		str += fmt.Sprintf("%d天", d)
	}
	if h > 0 {
		str += fmt.Sprintf("%d小时", h)
	}
	if m > 0 {
		str += fmt.Sprintf("%d分", m)
	}
	if s > 0 || str == "" {
		str += fmt.Sprintf("%d秒", s)
	}
	return str
}

func FormatDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())

}

//获取昨日0点时间戳
func GetYesterdayZeroTimestamp() int {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return int(tm1.Unix()) - 86400
}

// 获取该日0点时间戳
func GetThatZeroTimestamp(timestamp int64) int {
	t := time.Unix(timestamp, 0)
	t1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return int(t1.Unix())
}

// 获取当前时间戳
func GetTimestampMicro() int64 {
	return gtime.TimestampMilli()
}

// 获取当前时间戳
func GetTimestamp() int {
	return int(time.Now().Unix())
}

//获取相差时间
func GetTimeDifferS(startTime, endTime string) (second int) {
	t1 := gtime.ParseTimeFromContent(startTime, "2006-01-02 15:04:05")
	t2 := gtime.ParseTimeFromContent(endTime, "2006-01-02 15:04:05")
	//t1, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	//t2, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if t1 != nil && t2 != nil && t1.Before(t2) {
		second = gconv.Int(t2.Unix() - t1.Unix())
		return second
	} else {
		return second
	}
}

func TimeStrToStamp(timeStr string) int {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	CheckError(err)
	return int(t.Unix())
}
