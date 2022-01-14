package xifs

// 路由抽象接口
// 路由里的数据都是XRequest
type XRouter interface {

	// 处理之前的方法
	PreHandle(request XRequest)
	// 处理主方法
	Handle(request XRequest)
	// 处理之后的方法
	PostHandle(request XRequest)
}
