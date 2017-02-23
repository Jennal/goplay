package aop

import (
	// "fmt"
	"time"
)

type Aspect interface {
	Recover(err func(interface{})) Aspect
	Retry(times int) Aspect
	Repeat(times int) Aspect
	Delay(d time.Duration) Aspect
	Do(work func())
}

func NewAspect() Aspect {
	return &aspect{func(do func()) {
		do()
	}}
}

type aspect struct {
	chain func(func())
}

func (a *aspect) Join(work func(func())) Aspect {
	// fmt.Println("Joining", work)
	old_chain := a.chain
	a.chain = func(do func()) {
		// fmt.Println("chaining", work)
		old_chain(func() {
			work(do)
		})
	}

	return a
}

func (a *aspect) Recover(onErr func(interface{})) Aspect {
	a.Join(func(do func()) {
		Recover(do, onErr)
	})

	return a
}

func (a *aspect) Retry(times int) Aspect {
	// fmt.Println("Joining Retry", times)
	a.Join(func(do func()) {
		Retry(do, times)
	})

	return a
}

func (a *aspect) Repeat(times int) Aspect {
	// fmt.Println("Joining Repeat", times)
	a.Join(func(do func()) {
		Repeat(do, times)
	})

	return a
}

func (a *aspect) Delay(d time.Duration) Aspect {
	a.Join(func(do func()) {
		Delay(do, d)
	})

	return a
}

func (a *aspect) Do(work func()) {
	// fmt.Println("Doing", work)
	if a.chain != nil {
		// fmt.Printf("chain %#v\n", a.chain)
		a.chain(work)
	} else {
		// fmt.Println("work", work)
		work()
	}
}
