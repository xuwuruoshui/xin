package main

import (
	"log"
	"xin/xifs"
	"xin/xnet"
)

// 实现XRouter时，可以先实现这个BaseRouter，然后根据自己的需求修改
type PingRouter struct {
}

// 处理之前的方法
func (p *PingRouter) PreHandle(req xifs.XRequest) {
	log.Println("Call Router PreHandle")
	_, err := req.Connection().TCPConnetion().Write([]byte("before ping...\n"))
	if err != nil {
		log.Println("call back before ping error")
	}
}

// 处理主方法
func (p *PingRouter) Handle(req xifs.XRequest) {
	log.Println("Call Router Handle")
	_, err := req.Connection().TCPConnetion().Write([]byte("ping...\n"))
	if err != nil {
		log.Println("call back ping error")
	}
}

// 处理之后的方法
func (p *PingRouter) PostHandle(req xifs.XRequest) {
	log.Println("Call Router PostHandle")
	_, err := req.Connection().TCPConnetion().Write([]byte("post ping...\n"))
	if err != nil {
		log.Println("call back post ping error")
	}
}

func main() {
	// 1、创建xin的server
	s := xnet.NewServer()

	// 2、给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 2、启动xin的server
	s.Run()
}
