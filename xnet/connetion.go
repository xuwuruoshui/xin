package xnet

import (
	"errors"
	"github.com/xuwuruoshui/xin/config"
	"github.com/xuwuruoshui/xin/xifs"
	"log"
	"net"
)

type Connection struct {
	// 当前连接的socket: TCP套接字
	conn *net.TCPConn

	// 链接的ID
	connID uint32

	// 当前的链接状态
	isClosed bool

	// 告知当前连接已经退出的/停止 channel
	exitChan chan bool

	// 用于Read、Write的channel(无缓冲)
	msgChan chan []byte

	// 该链接处理的方法Router
	msgHandler xifs.XMessageHandle
}

// 初始化链接模块的方法
func NewConnetion(conn *net.TCPConn, connID uint32, msgHandler xifs.XMessageHandle) *Connection {
	return &Connection{
		conn:       conn,
		connID:     connID,
		msgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		exitChan:   make(chan bool, 1),
	}
}

// 读业务
func (c *Connection) StartReader() {
	log.Println("[Reader Goroutine is running]...")
	defer log.Printf("[Reader is exit,remote addr is %s],connID=%d\n", c.RemoteAddr().String(), c.connID)
	defer c.Stop()

	for {
		// 创建一个datapack对象
		datapack := NewDataPack()
		headData := make([]byte, datapack.GetHeadLen())

		if _, err := c.conn.Read(headData); err != nil {
			log.Println("read headLen error:", err)
			break
		}

		msg, err := datapack.Unpack(headData)
		if err != nil {
			log.Println("Unpack data error:", err)
			break
		}

		if msg.GetLength() > 0 {
			msg.SetData(make([]byte, msg.GetLength()))
			if _, err = c.conn.Read(msg.GetData()); err != nil {
				log.Println("read Data error:", err)
				break
			}
		}

		req := Request{
			conn: c,
			msg:  msg,
		}

		if config.GloabalConf.WorkerPoolSize > 0 {
			// 已经开启了工作池，将消息发送给工作池
			c.msgHandler.SendMsgToTaskQueue(&req)

		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			go c.msgHandler.DoMsgHandle(&req)
		}

	}
}

// 写业务
func (c *Connection) StartWrite() {
	log.Println("[Writer Gorutine is running]...")
	defer log.Printf("%s Writer exit!!!", c.RemoteAddr())

	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.conn.Write(data); err != nil {
				log.Printf("Send data error%s,conn writer exit!!!", err)
				return
			}
		case <-c.exitChan:
			// Reader结束,Writer也退出
			return

		}
	}
}

// 发送数据到远程客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection is closed when send msg !!!")
	}
	// 封包
	msg := NewMessage(msgId, data)
	datapack := NewDataPack()
	packData, err := datapack.Pack(msg)

	if err != nil {
		log.Println("pack data error:", err)
		return err
	}

	// 数据发送到channel
	c.msgChan <- packData
	return nil
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	log.Println("conn Start() ... connID=", c.connID)
	// 启动从当前链接
	// TODO 启动从当前链接写数据的业务
	go c.StartReader()
	go c.StartWrite()
}

// 停止链接	结束当前连接的工作
func (c *Connection) Stop() {
	log.Println("conn Stop().. connID=", c.connID)

	// 如果当前连接已经关闭
	if c.isClosed {
		return
	}

	c.isClosed = true
	// 关闭链接
	c.conn.Close()

	// 告知writer关闭
	c.exitChan <- true

	// 回收资源
	close(c.exitChan)
	close(c.msgChan)
}

// 获取当前链接的绑定 socket conn
func (c *Connection) TCPConnetion() *net.TCPConn {
	return c.conn
}

// 获取当前链接模块的链接ID
func (c *Connection) ConnId() uint32 {
	return c.connID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
