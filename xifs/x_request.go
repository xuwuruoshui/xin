package xifs

type XRequest interface {

	// 得到当前链接
	Connection() XConnection
	// 得到当前消息数据
	Data() []byte
}
