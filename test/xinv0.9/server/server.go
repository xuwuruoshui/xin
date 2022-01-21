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

	if err := req.Connection().SendMsg(req.MsgId(), []byte("XinV0.9 server ping...")); err != nil {
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

	if err := req.Connection().SendMsg(req.MsgId(), []byte("XinV0.9 server Hello...")); err != nil {
		log.Println("send msg error:", err)
	}
}

// 创建链接之后的钩子函数
func afterConnection(conn xifs.XConnection) {
	log.Println("BeforeConnection is Called...")
	err := conn.SendMsg(111, []byte("BeforeConnection Begin"))
	if err != nil {
		log.Println(err)
	}
}

// 销毁链接之前的钩子函数
func beforeStop(conn xifs.XConnection) {
	log.Println("BeforeStop is Called...")
	log.Printf("conn Id=%d is Lost...", conn.ConnId())
}

func main() {
	// 1、创建xin的server
	s := xnet.NewServer()

	// 2、注册链接Hook钩子函数
	s.SetOnConnStart(afterConnection)
	s.SetOnConnStop(beforeStop)

	// 3、给当前zinx框架添加一个自定义的router
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	// 3、启动xin的server
	s.Run()
}
