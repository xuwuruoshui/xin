package xnet

import (
	"fmt"
	"log"
	"net"
	"time"
	"xin/xifs"
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
}

// 启动
func (s *Server) Start() {
	fmt.Printf("[Start]Server Listener at IP:%s,Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 以下的API更底层一点
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

		fmt.Printf("Start xin server success: %s is Listening", s.Name)

		// 3、阻塞等待客户端连接
		for {
			// 如果有客户端过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 已经与客户端建立连接，v1.0就实现一个512k回显业务
			go func() {
				for {
					buf := make([]byte, 512)
					c, err := conn.Read(buf)
					if err != nil {
						log.Println("Read err:", err)
						time.Sleep(time.Second)
						continue
					}
					log.Println(string(buf[:c]))

					// 回显
					if _, err := conn.Write(buf[:c]); err != nil {
						log.Println("Write err:", err)
						continue
					}
				}
			}()
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

func NewServer(name string) xifs.XServer {
	return &Server{Name: name, IPVersion: "tcp4", IP: "0.0.0.0", Port: 9999}
}
