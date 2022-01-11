package xnet

import (
	"log"
	"net"
	"xin/xifs"
)

type Connection struct {
	// 当前连接的socket: TCP套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 当前连接所绑定的处理业务方法API
	handleAPI xifs.HandleFunc

	// 告知当前连接已经退出的/停止 channel
	ExitChan chan bool
}

// 初始化链接模块的方法
func NewConnetion(conn *net.TCPConn,connID uint32,callback_api xifs.HandleFunc)  *Connection{
	return &Connection{
		Conn: conn,
		ConnID: connID,
		handleAPI: callback_api,
		isClosed: false,
		ExitChan: make(chan bool,1),
	}
}

// 读业务
func (c *Connection)StartReader(){
	log.Println("Reader Goroutine is running...")
	defer log.Println("connID=",c.ConnID,"Reader is exit,remote addr is",c.RemoteAddr().String())
	defer c.Stop()

	for  {
		// 读取客户端的数据到buf中，最大512字节
		data := make([]byte,512)
		n,err := c.Conn.Read(data)
		if err!=nil{
			log.Println("recv data err:",err)
			continue
		}
		log.Println(string(data[:n]))

		// 调用当前连接所绑定的HandleAPI
		if err:= c.handleAPI(c.Conn,data,n);err!=nil{
			log.Println("ConnID",c.ConnID,"handle is error",err)
			break
		}

	}
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection)Start(){
	log.Println("Conn Start() ... ConnID=",c.ConnID)
	// 启动从当前链接
	// TODO 启动从当前链接写数据的业务
	c.StartReader()
}

// 停止链接	结束当前连接的工作
func (c *Connection)Stop(){
	log.Println("Conn Stop().. ConnID=",c.ConnID)

	// 如果当前连接已经关闭
	if c.isClosed{
		return
	}

	c.isClosed = true
	// 关闭链接
	c.Conn.Close()

	// 回收资源
	close(c.ExitChan)
}

// 获取当前链接的绑定 socket conn
func (c *Connection)TCPConnetion() *net.TCPConn{
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection)ConnId() uint32{
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection)RemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据发送给远程的客户端
func (c *Connection)Send(data []byte)error{
	return nil
}