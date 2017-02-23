package aop

import (
	"fmt"
	"testing"
	"time"
)

func TestAop(t *testing.T) {
	NewAspect().
		Retry(3).
		Delay(10 * time.Second).
		Repeat(5).
		Do(func() {
		fmt.Println("Test")
	})
}
