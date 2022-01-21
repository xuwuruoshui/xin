package xifs

// 链接管理模块抽象层
type XConnectionManager interface {
	// 添加链接
	Add(conn XConnection)
	// 删除链接
	Remove(conn XConnection)
	//根据connID获取链接
	Get(connId uint32) (XConnection, error)
	// 得到当前连接总数
	Len() uint
	//清除并终止所有的连接
	ClearConn()
}
