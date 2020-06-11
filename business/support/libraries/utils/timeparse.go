package utils

import (
	"business/support/libraries/loggers"
	"fmt"
	"time"
)

var local *time.Location

func init() {
	var err error
	local, err = time.LoadLocation("Asia/Chongqing") //等同于"UTC"
	if err != nil {
		loggers.Error.Printf("ConvertTimestamp2CHDateStr LoadLocation 'Asia/Chongqing' failed!,error:%v", err)
		panic(err)
	}
}

// 将时间戳转换成北京时区的日期
func Timestamp2CHDateStr(timestamp int64) (string, error) {
	ti := time.Unix(timestamp, 0)
	formate := "2006-01-02"
	return fmt.Sprintf(ti.In(local).Format(formate)), nil
}

// 将时间戳转换成北京时区的时间,最小单位是分钟
func Timestamp2CHMinuteStr(timestamp int64) (string, error) {
	ti := time.Unix(timestamp, 0)
	formate := "2006-01-02 15:04:00"
	return fmt.Sprintf(ti.In(local).Format(formate)), nil
}

func Timestamp2MinuteTimestamp(timestamp int64) int64 {
	return Timestamp2Unit(timestamp, 60)
}

func TimestampTo5MinuteTimestamp(timestamp int64) int64 {
	return Timestamp2Unit(timestamp, 5*60)
}

func Timestamp2HourTimestamp(timestamp int64) int64 {
	return Timestamp2Unit(timestamp, 60*60)
}

func Timestamp2HalfHourTimestamp(timestamp int64) int64 {
	return Timestamp2Unit(timestamp, 30*60)
}

func Timestamp2DayTimestamp(timestamp int64) int64 {
	return Timestamp2Unit(timestamp, 24*60*60)
}

func Timestamp2Unit(timestamp int64, unit int64) int64 {
	// 由于时间采用的是单位进制方式,所以在时间取完摩之后,需要判断是否属于下一个进制
	// 如果属于下一个进制需要加1
	ts := timestamp - (timestamp % unit)
	if timestamp > ts {
		ts = ts + 1
	} else {
		ts = ts - unit + 1
		if ts < 1 {
			ts = 1
		}
	}
	return ts
}

func TodayTimestamp() int64 {
	return Timestamp2DayTimestamp(time.Now().UTC().Unix())
}

func NowMinutes() int64 {
	timestamp := time.Now().UTC().Unix()
	return timestamp - timestamp%60
}

func NowHalfHours() int64 {
	timestamp := time.Now().UTC().Unix()
	return TimestampHalfHours(timestamp)
}

func TimestampHalfHours(timestamp int64) int64 {
	timestamp = timestamp - (timestamp % (1800))
	return timestamp / 1800
}

func TimestampHours(timestamp int64) int64 {
	timestamp = timestamp - (timestamp % (3600))
	return timestamp / 3600
}

func IsSameHalfHour(t1 int64, t2 int64) bool {
	t1 = t1 - (t1 % (1800))
	t2 = t2 - (t2 % (1800))
	if t1 == t2 {
		return true
	}
	return false
}

func Time2CHDateStr(ti time.Time) (string, error) {
	formate := "2006-01-02"
	return fmt.Sprintf(ti.In(local).Format(formate)), nil
}

func Time2CHsecondStr(ti time.Time) string {
	formate := "2006-01-02 15:04:05"
	return fmt.Sprintf(ti.In(local).Format(formate))
}

func CompareDateStr(date1 string, date2 string) (int, error) {
	t1, err := time.Parse("2006-01-02", date1)
	if nil != err {
		return -2, err
	}
	t2, err := time.Parse("2006-01-02", date2)
	if nil != err {
		return -2, err
	}
	if t1.Before(t2) {
		// 判断时间大小后的逻辑处理
		return -1, nil
	} else if t1.Equal(t2) {
		return 0, nil
	} else {
		return 1, nil
	}
}

func IsSameDate(timestamp1 int64, timestamp2 int64) bool {
	date1, err := Timestamp2CHDateStr(timestamp1)
	if nil != err {
		return false
	}
	date2, err := Timestamp2CHDateStr(timestamp2)
	if nil != err {
		return false
	}
	comp, err := CompareDateStr(date1, date2)
	if nil != err {
		return false
	}
	if 0 == comp {
		return true
	}
	return false
}

func DisDays(timestamp1 int64, timestamp2 int64) (int, error) {
	date1, err := Timestamp2CHDateStr(timestamp1)
	if nil != err {
		return -111111, err
	}
	date2, err := Timestamp2CHDateStr(timestamp2)
	if nil != err {
		return -111111, err
	}
	t1, err := time.Parse("2006-01-02", date1)
	if nil != err {
		return -111111, err
	}
	t2, err := time.Parse("2006-01-02", date2)
	if nil != err {
		return -111111, err
	}
	diff := t2.Unix() - t1.Unix() //
	days := int(diff / 3600 / 24)
	return days, nil
}

func DisDays2(timestamp1 time.Time, timestamp2 time.Time) (int, error) {
	formate := "2006-01-02"
	date1 := fmt.Sprintf(timestamp1.In(local).Format(formate))
	date2 := fmt.Sprintf(timestamp2.In(local).Format(formate))
	t1, err := time.Parse("2006-01-02", date1)
	if nil != err {
		return -111111, err
	}
	t2, err := time.Parse("2006-01-02", date2)
	if nil != err {
		return -111111, err
	}
	diff := t2.Unix() - t1.Unix() //
	days := int(diff / 3600 / 24)
	return days, nil
}

func SimpleDisDays(timestamp1 int64, timestamp2 int64) int {
	diff := timestamp1 - timestamp2 //
	days := int(diff / 3600 / 24)
	return days
}

func DisUTCDays(timestamp1 int64, timestamp2 int64) int {
	t1 := time.Unix(timestamp1, 0)
	t2 := time.Unix(timestamp2, 0)
	return t1.Day() - t2.Day()
}

func NowBeijingTime() time.Time {
	ti := time.Now()
	return ti.In(local)
}

func ToBeijingTime(utcTime *time.Time) time.Time {
	return utcTime.In(local)
}
