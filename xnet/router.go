package xnet

import "xin/xifs"

// 实现XRouter时，可以先实现这个BaseRouter，然后根据自己的需求修改
type BaseRouter struct {
}

// 处理之前的方法
func (br *BaseRouter) PreHandle(request xifs.XRequest) {}

// 处理主方法
func (br *BaseRouter) Handle(request xifs.XRequest) {}

// 处理之后的方法
func (br *BaseRouter) PostHandle(request xifs.XRequest) {}
