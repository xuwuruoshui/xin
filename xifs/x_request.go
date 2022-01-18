package xifs

type XRequest interface {

	// 得到当前链接
	Connection() XConnection
	// 得到当前消息数据
	Data() []byte
	// 得到当前消息ID
	MsgId() uint32
}
