package event

import (
	"fmt"
	"sync"
	"testing"
)

func TestEvent(t *testing.T) {
	// 事件注册
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		Register(EVENT_A, func(args ...any) {
			a := args[0].(int)
			b := args[1].(int)
			res := a + b
			fmt.Println(res)
		})
		wg.Done()
	}()
	go func() {
		Register(EVENT_A, func(args ...any) {
			a := args[0].(int)
			b := args[1].(int)
			res := a * b
			fmt.Println(res)
		})
		wg.Done()
	}()
	wg.Wait()

	// 事件触发
	wg.Add(2)
	go func() {
		Fire(EVENT_A, 1, 2)
		wg.Done()
	}()
	go func() {
		Fire(EVENT_A, 3, 4)
		wg.Done()
	}()
	wg.Wait()
}
