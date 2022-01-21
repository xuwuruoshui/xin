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
	connections map[uint32]xifs.XConnection
	// 链接读写锁
	connectionLock sync.RWMutex
}

// 创建当前连接的方法
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]xifs.XConnection),
	}
}

// 添加链接
func (c *ConnectionManager) Add(conn xifs.XConnection) {
	// 保护共享资源map，加写锁
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	// 将conn加入到ConnManager
	c.connections[conn.ConnId()] = conn
	log.Printf("connId=%d connection add to ConnManager successfully:conn num=%d\n", conn.ConnId(), c.Len())

}

// 删除链接
func (c *ConnectionManager) Remove(conn xifs.XConnection) {
	// 保护共享资源map，加写锁
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	// 删除连接信息
	delete(c.connections, conn.ConnId())

	log.Printf("connId=%d removed from  ConnManager successfully:conn num=%d\n", conn.ConnId(), c.Len())
}

//根据connID获取链接
func (c *ConnectionManager) Get(connId uint32) (xifs.XConnection, error) {
	c.connectionLock.RLock()
	defer c.connectionLock.RUnlock()

	if conn, ok := c.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found!!!")
	}
}

// 得到当前连接总数
func (c *ConnectionManager) Len() uint {
	return uint(len(c.connections))
}

//清除并终止所有的连接
func (c *ConnectionManager) ClearConn() {
	// 保存共享资源map，加写锁
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	// 删除conn并停止conn的工作
	for connId, conn := range c.connections {
		// 停止
		conn.Stop()
		delete(c.connections, connId)
	}
	log.Printf("Clear All connection success!!! Conn num=%d\n", c.Len())
}
