package goutil

import "time"

func Millisecond() int64 {
	unixNano := time.Now().UnixNano()
	umillisec := unixNano / 1000000
	return umillisec
}
