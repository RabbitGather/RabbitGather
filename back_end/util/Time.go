package util

import "time"

//func UnixTimeAfterSec(sec time.Duration) int64 {
//	nowTime := time.Now()
//	return nowTime.Add(sec * time.Second).Unix()
//}

//func UnixNow() int64 {
//	return time.Now().Unix()
//}

func RunAfterFuncWithCancel(timeout time.Duration, f func(), cancel <-chan struct{}) {
	go func() {
		for true {
			select {
			case <-cancel:
				return
			case <-time.After(timeout):
				f()
			}
		}
	}()

}
