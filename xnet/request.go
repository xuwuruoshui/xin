package xnet

import "github.com/xuwuruoshui/xin/xifs"

type Request struct {
	// 已经和客户端建立好的链接
	Conn xifs.XConnection

	// 客户端请求的数据
	Msg xifs.XMessage
}

// 获取链接
func (r *Request) Connection() xifs.XConnection {
	return r.Conn
}

// 获取请求数据
func (r *Request) Data() []byte {
	return r.Msg.GetData()
}

// 得到当前消息ID
func (r *Request) MsgId() uint32 {
	return r.Msg.GetId()
}
