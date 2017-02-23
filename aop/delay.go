package aop

import (
	"time"
)

func Delay(work func(), d time.Duration) {
	time.Sleep(d)
	work()
}
