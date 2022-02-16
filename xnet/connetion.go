package xnet

import (
	"errors"
	"github.com/xuwuruoshui/xin/config"
	"github.com/xuwuruoshui/xin/xifs"
	"log"
	"net"
	"sync"
)

type Connection struct {

	// 当前Conn与哪个server相关联
	TcpServer xifs.XServer

	// 当前连接的socket: TCP套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前的链接状态
	IsClosed bool

	// 告知当前连接已经退出的/停止 channel
	ExitChan chan bool

	// 用于Read、Write的channel(无缓冲)
	MsgChan chan []byte

	// 该链接处理的方法Router
	MsgHandler xifs.XMessageHandle

	// 链接属性集合
	Property map[string]interface{}

	// 保护链接属性的锁
	PropertyLock sync.RWMutex
}

// 初始化链接模块的方法
func NewConnetion(server xifs.XServer, conn *net.TCPConn, connID uint32, msgHandler xifs.XMessageHandle) *Connection {

	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		IsClosed:   false,
		MsgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		Property:   make(map[string]interface{}),
	}
	// 将当前链接添加到链接管理器ConnectionManager
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// 读业务
func (c *Connection) StartReader() {
	log.Println("[Reader Goroutine is running]...")
	defer log.Printf("[Reader is exit,remote addr is %s],ConnID=%d\n", c.RemoteAddr().String(), c.ConnID)
	defer c.Stop()

	for {
		// 创建一个datapack对象
		datapack := NewDataPack()
		headData := make([]byte, datapack.GetHeadLen())

		if _, err := c.Conn.Read(headData); err != nil {
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
			if _, err = c.Conn.Read(msg.GetData()); err != nil {
				log.Println("read Data error:", err)
				break
			}
		}

		req := Request{
			Conn: c,
			Msg:  msg,
		}

		if config.GloabalConf.WorkerPoolSize > 0 {
			// 已经开启了工作池，将消息发送给工作池
			c.MsgHandler.SendMsgToTaskQueue(&req)

		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			go c.MsgHandler.DoMsgHandle(&req)
		}

	}
}

// 写业务
func (c *Connection) StartWrite() {
	log.Println("[Writer Gorutine is running]...")
	defer log.Printf("%s Writer exit!!!", c.RemoteAddr())

	for {
		select {
		case data := <-c.MsgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				log.Printf("Send data error%s,Conn writer exit!!!", err)
				return
			}
		case <-c.ExitChan:
			// Reader结束,Writer也退出
			return

		}
	}
}

// 发送数据到远程客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("Connection is closed when send Msg !!!")
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
	c.MsgChan <- packData
	return nil
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	log.Println("Conn Start() ... ConnID=", c.ConnID)
	// 启动从当前链接
	// TODO 启动从当前链接写数据的业务
	go c.StartReader()
	go c.StartWrite()

	// 开发者传递进来，创建链接之后需要调用处理的业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止链接	结束当前连接的工作
func (c *Connection) Stop() {
	log.Println("Conn Stop().. ConnID=", c.ConnID)

	// 如果当前连接已经关闭
	if c.IsClosed {
		return
	}

	c.IsClosed = true

	// 调用开发者注册的Hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭链接
	c.Conn.Close()

	// 告知writer关闭
	c.ExitChan <- true

	// 将当前链接从Connection Manager中删除
	c.TcpServer.GetConnMgr().Remove(c)
	// 回收资源
	close(c.ExitChan)
	close(c.MsgChan)
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()

	// 添加一个链接属性
	c.Property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.PropertyLock.RLock()
	defer c.PropertyLock.RUnlock()

	// 读取属性
	if value, ok := c.Property[key]; ok {
		return value, nil
	}
	return nil, errors.New("Property not found!!!")
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()
	delete(c.Property, key)
}
