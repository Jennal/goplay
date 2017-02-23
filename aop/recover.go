package aop

func Recover(work func(), onErr func(interface{})) {
	defer func() {
		err := recover()
		if err != nil {
			onErr(err)
		}
	}()

	work()
}
