package xnet

import (
	"errors"
	"github.com/xuwuruoshui/xin/xifs"
	"log"
	"sync"
)

// 链接管理
type ConnectionManager struct {
	// 链接map
	Connections map[uint32]xifs.XConnection
	// 链接读写锁
	ConnectionLock sync.RWMutex
}

// 创建当前连接的方法
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Connections: make(map[uint32]xifs.XConnection),
	}
}

// 添加链接
func (c *ConnectionManager) Add(conn xifs.XConnection) {
	// 保护共享资源map，加写锁
	c.ConnectionLock.Lock()
	defer c.ConnectionLock.Unlock()

	// 将conn加入到ConnManager
	c.Connections[conn.GetConnId()] = conn
	log.Printf("connId=%d connection add to ConnManager successfully:Conn num=%d\n", conn.GetConnId(), c.Len())

}

// 删除链接
func (c *ConnectionManager) Remove(conn xifs.XConnection) {
	// 保护共享资源map，加写锁
	c.ConnectionLock.Lock()
	defer c.ConnectionLock.Unlock()

	// 删除连接信息
	delete(c.Connections, conn.GetConnId())

	log.Printf("connId=%d removed from  ConnManager successfully:Conn num=%d\n", conn.GetConnId(), c.Len())
}

//根据connID获取链接
func (c *ConnectionManager) Get(connId uint32) (xifs.XConnection, error) {
	c.ConnectionLock.RLock()
	defer c.ConnectionLock.RUnlock()

	if conn, ok := c.Connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found!!!")
	}
}

// 得到当前连接总数
func (c *ConnectionManager) Len() uint {
	return uint(len(c.Connections))
}

//清除并终止所有的连接
func (c *ConnectionManager) ClearConn() {
	// 保存共享资源map，加写锁
	c.ConnectionLock.Lock()
	defer c.ConnectionLock.Unlock()

	// 删除conn并停止conn的工作
	for connId, conn := range c.Connections {
		// 停止
		conn.Stop()
		delete(c.Connections, connId)
	}
	log.Printf("Clear All connection success!!! Conn num=%d\n", c.Len())
}
