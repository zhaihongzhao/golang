package chanrpc

import (
    "testing"
    "time"
)

func TestChanRPC(t *testing.T) {
    // 服务器
    s := NewServer(100)
    s.Register("add", func(args []any) []any {
        a := args[0].(int)
        b := args[1].(int)
        res := a + b
        return []any{res}
    })
    s.Register("mult", func(args []any) []any {
        a := args[0].(int)
        b := args[1].(int)
        res := a * b
        return []any{res}
    })
    s.Start()

    // 客户端1
    go func() {
        c := NewClient(10)
        c.Attach(s)
        // 同步模式
        ret, err := c.SyncCall("add", 1, 2)
        if err != nil {
            t.Error(err)
        } else {
            t.Log(ret[0])
        }
        ret, err = c.SyncCall("mult", 1, 2)
        if err != nil {
            t.Error(err)
        } else {
            t.Log(ret[0])
        }
        // 异步模式
        cb := func(ret []any, err error) {
            if err != nil {
                t.Error(err)
            } else {
                t.Log(ret[0])
            }
        }
        c.AsynCall("add", cb, 1, 2)
        c.AsynCall("mult", cb, 1, 2)
        // Go模式
        c.Go("add", 1, 2)
        c.Go("mult", 1, 2)
    }()

    // 客户端2
    go func() {
        c := NewClient(10)
        c.Attach(s)
        // 同步模式
        ret, err := c.SyncCall("add", 3, 4)
        if err != nil {
            t.Error(err)
        } else {
            t.Log(ret[0])
        }
        ret, err = c.SyncCall("mult", 3, 4)
        if err != nil {
            t.Error(err)
        } else {
            t.Log(ret[0])
        }
        // 异步模式
        cb := func(ret []any, err error) {
            if err != nil {
                t.Error(err)
            } else {
                t.Log(ret[0])
            }
        }
        c.AsynCall("add", cb, 3, 4)
        c.AsynCall("mult", cb, 3, 4)
        // Go模式
        c.Go("add", 3, 4)
        c.Go("mult", 3, 4)
    }()

    time.Sleep(time.Second * 2)
}
