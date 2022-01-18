package xifs

// 封包、拆包,解决粘包问题

type XDataPack interface {

	// 获取包的头的长度方法
	GetHeadLen() uint32

	// 封包方法
	Pack(msg XMessage) ([]byte, error)

	// 拆包方法
	Unpack([]byte) (XMessage, error)
}
