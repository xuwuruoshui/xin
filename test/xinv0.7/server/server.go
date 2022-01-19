package main

import (
	"github.com/xuwuruoshui/xin/xifs"
	"github.com/xuwuruoshui/xin/xnet"
	"log"
)

// 实现XRouter时，可以先实现这个BaseRouter，然后根据自己的需求修改
type PingRouter struct {
	xnet.BaseRouter
}

// 处理主方法
func (p *PingRouter) Handle(req xifs.XRequest) {
	log.Println("Call Router Handle")

	// 先读取客户端数据，再回写ping...
	log.Printf("Recv from clien:msgId=%d,data=%s\n", req.MsgId(), req.Data())

	if err := req.Connection().SendMsg(req.MsgId(), []byte("XinV0.7 server ping...")); err != nil {
		log.Println("send msg error:", err)
	}
}

type HelloRouter struct {
	xnet.BaseRouter
}

// 处理主方法
func (h *HelloRouter) Handle(req xifs.XRequest) {
	log.Println("Call Router Handle")

	// 先读取客户端数据，再回写ping...
	log.Printf("Recv from clien:msgId=%d,data=%s\n", req.MsgId(), req.Data())

	if err := req.Connection().SendMsg(req.MsgId(), []byte("XinV0.6 server Hello...")); err != nil {
		log.Println("send msg error:", err)
	}
}

func main() {
	// 1、创建xin的server
	s := xnet.NewServer()

	// 2、给当前zinx框架添加一个自定义的router
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	// 3、启动xin的server
	s.Run()
}
