package event

/*
事件系统的基本原理：
‌事件注册：将事件ID和响应函数关联并保存起来。
‌事件触发：将事件ID和事件发生的参数发送到事件系统。
‌事件响应：‌通过事件ID找到对应的响应函数，按注册顺序依次调用。
事件系统能有效地将事件触发和事件响应两端代码解耦。
*/

import (
	"sync"
)

type Handler func(...any)

var eventMap = make(map[int][]Handler)
var mu sync.Mutex

// 事件注册
func Register(id int, h Handler) {
	mu.Lock()
	defer mu.Unlock()
	list := eventMap[id]
	list = append(list, h)
	eventMap[id] = list
}

// 事件触发
func Fire(id int, args ...any) {
	list := eventMap[id]
	for _, h := range list {
		h(args...)
	}
}

// 事件定义
const (
	EVENT = iota
	EVENT_A
	EVENT_B
	EVENT_C
)
