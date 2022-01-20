package xnet

import (
	"fmt"
	"github.com/xuwuruoshui/xin/config"
	"github.com/xuwuruoshui/xin/xifs"
	"log"
	"net"
)

// 实现XServer
type Server struct {
	// 服务器名称
	Name string
	// 版本
	IPVersion string
	// IP地址
	IP string
	// 端口
	Port int
	// 当前的Server添加一个MessageHandler,用于msgId和router的绑定
	msgHandler xifs.XMessageHandle
}

// 启动
func (s *Server) Start() {
	log.Printf("[Xin]Server Listener at Name:%s,Version %s.\n", config.GloabalConf.Name, config.GloabalConf.Version)
	log.Printf("[Xin]Server Listener at IP:%s,Port %d, is starting\n", config.GloabalConf.Host, config.GloabalConf.Port)
	log.Printf("[Xin]Server Listener at MaxConn:%d,MaxPackageSize %d, is starting\n", config.GloabalConf.MaxConn, config.GloabalConf.MaxPackageSize)

	go func() {
		//0、开启消息队列及Worker工作池
		s.msgHandler.StartWorkerPool()

		// 1、获取一个TCP的Address
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		// 2、监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen:", addr, ",err:", err)
			return
		}

		fmt.Printf("Start xin server success: %s is Listening\n", s.Name)
		var cid uint32 = 0

		// 3、阻塞等待客户端连接
		for {
			// 如果有客户端过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnetion(conn, cid, s.msgHandler)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()

}

// 停止
func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或一些开辟的链接信息进行停止或者回收
}

// 运行
func (s *Server) Run() {
	s.Start()

	//TODO 做一些启动服务器之外的额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgId uint32, router xifs.XRouter) {
	s.msgHandler.AddRouter(msgId, router)
	log.Println("AddRouter Success!!!")
}

func NewServer() xifs.XServer {
	return &Server{
		Name:       config.GloabalConf.Name,
		IPVersion:  "tcp4",
		IP:         config.GloabalConf.Host,
		Port:       config.GloabalConf.Port,
		msgHandler: NewMsgHandle()}
}
