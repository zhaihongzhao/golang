package chanrpc

/*
提供一套基于channel的RPC机制，用于服务器各模块之间的通信。
chanRPC支持三种模式：
‌同步模式‌：调用并等待RPC返回结果。
‌异步模式‌：调用并提供回调函数，回调函数在RPC返回后被调用。
‌Go模式‌：调用并立即返回，忽略任何返回值和错误‌。
*/

import (
	"errors"
	"fmt"
	"time"
)

type Func func([]any) []any
type Cb func([]any, error)

// 服务器
type Server struct {
	mapFunc  map[string]Func // 函数表
	chanCall chan *CallInfo  // 调用信息
}

// 客户端
type Client struct {
	s           *Server       // 绑定的服务器
	chanSyncRet chan *RetInfo // 同步调用返回信息
	chanAsynRet chan *RetInfo // 异步调用返回信息
	asynCallNum int           // 进行中的异步调用数量
}

// 调用信息
type CallInfo struct {
	id      string        // 函数id
	args    []any         // 参数列表
	chanRet chan *RetInfo // 返回信息
	cb      Cb            // 回调函数
}

// 返回信息
type RetInfo struct {
	cb  Cb    // 回调函数
	ret []any // 返回值
	err error // 错误信息
}

// 新建服务器
func NewServer(size int) *Server {
	return &Server{
		mapFunc:  make(map[string]Func),
		chanCall: make(chan *CallInfo, size),
	}
}

// 启动服务器
func (s *Server) Start() {
	go func() {
		for ci := range s.chanCall {
			s.exec(ci)
		}
	}()
}

// 注册函数
func (s *Server) Register(id string, f Func) {
	if s.mapFunc[id] != nil {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}
	s.mapFunc[id] = f
}

// 执行函数
func (s *Server) exec(ci *CallInfo) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("%v", r)
			s.ret(ci, &RetInfo{err: err})
		}
	}()

	f := s.mapFunc[ci.id]
	if f == nil {
		panic(fmt.Sprintf("function id %v: not found", ci.id))
	}

	ret := f(ci.args)
	s.ret(ci, &RetInfo{ret: ret})
}

// 返回
func (s *Server) ret(ci *CallInfo, ri *RetInfo) {
	if ci.chanRet == nil {
		return
	}
	ri.cb = ci.cb
	ci.chanRet <- ri
	return
}

// 新建客户端
func NewClient(size int) *Client {
	c := &Client{
		chanSyncRet: make(chan *RetInfo, 1),
		chanAsynRet: make(chan *RetInfo, size),
	}

	go func() {
		// 读取异步调用的返回信息，调用回调函数
		for ri := range c.chanAsynRet {
			c.asynCallNum--
			ri.cb(ri.ret, ri.err)
		}
	}()

	return c
}

// 绑定服务器
func (c *Client) Attach(s *Server) {
	c.s = s
}

// 同步调用
func (c *Client) SyncCall(id string, args ...any) (ret []any, err error) {
	ci := &CallInfo{
		id:      id,
		args:    args,
		chanRet: c.chanSyncRet,
	}
	c.s.chanCall <- ci

	select {
	case ri := <-c.chanSyncRet:
		ret, err = ri.ret, ri.err
	case <-time.After(time.Second):
		ret, err = nil, errors.New("Timeout!")
	}
	return
}

// 异步调用
func (c *Client) AsynCall(id string, cb Cb, args ...any) {
	if c.asynCallNum >= cap(c.chanAsynRet) {
		cb(nil, errors.New("too many calls"))
		return
	}

	c.asynCallNum++
	ci := &CallInfo{
		id:      id,
		args:    args,
		chanRet: c.chanAsynRet,
		cb:      cb,
	}

	select {
	case c.s.chanCall <- ci:
	default:
		err := errors.New("chanrpc channel full")
		c.chanAsynRet <- &RetInfo{cb: cb, err: err}
	}
}

// Go模式
func (c *Client) Go(id string, args ...any) {
	c.s.chanCall <- &CallInfo{id: id, args: args}
}
