package test

import (
	//"Mockfish/engine"
	"time"
)

func getTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// func perftTest() {
// 	start_time := getTimestamp()

// }