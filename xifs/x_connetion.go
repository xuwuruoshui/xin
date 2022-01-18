package xifs

import "net"

type XConnection interface {

	// 启动链接 让当前的链接准备开始工作
	Start()
	// 停止链接	结束当前连接的工作
	Stop()
	// 获取当前链接的绑定 socket conn
	TCPConnetion() *net.TCPConn
	// 获取当前链接模块的链接ID
	ConnId() uint32
	// 获取远程客户端的TCP状态 IP port
	RemoteAddr() net.Addr
	// 发送数据，将数据发送给远程的客户端
	SendMsg(msgId uint32, data []byte) error
}
