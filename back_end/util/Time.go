package util

import "time"

//func UnixTimeAfterSec(sec time.Duration) int64 {
//	nowTime := time.Now()
//	return nowTime.Add(sec * time.Second).Unix()
//}

func UnixNow() int64 {
	return time.Now().Unix()
}
