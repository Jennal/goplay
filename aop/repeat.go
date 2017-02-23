package aop

func Repeat(work func(), times int) {
	for i := 0; i < times; i++ {
		work()
	}
}
