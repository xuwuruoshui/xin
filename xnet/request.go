package xnet

import "github.com/xuwuruoshui/xin/xifs"

type Request struct {
	// 已经和客户端建立好的链接
	conn xifs.XConnection

	// 客户端请求的数据
	data []byte
}

// 获取链接
func (r *Request) Connection() xifs.XConnection {
	return r.conn
}

// 获取请求数据
func (r *Request) Data() []byte {
	return r.data
}
