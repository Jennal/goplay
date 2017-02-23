package aop

import (
// "fmt"
)

func Retry(work func(), times int) {
	defer func() {
		err := recover()
		if err == nil || times <= 1 {
			return
		}

		times--
		Retry(work, times)
	}()

	// fmt.Println("Retring", times)
	work()
}
