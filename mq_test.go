package messagequeue

import (
    "fmt"
    "testing"
    "time"
)

type myMessage struct {
    info string
}

func (m *myMessage) Do() {
    fmt.Printf("Consumed: %s\n", m.info)
}

func TestMessageQueue(t *testing.T) {
    mq := New(10)
    mq.Start()
    defer mq.Stop()

    // 添加消息
    go func() {
        for i := 1; i <= 100; i++ {
            msg := &myMessage{info: fmt.Sprintf("job %d", i)}
            mq.Add(msg)
        }
    }()

    // 等待所有工作完成
    time.Sleep(time.Second * 2)
}
