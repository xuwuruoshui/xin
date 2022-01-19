package xifs

// 消息管理
type XMessageHandle interface {

	// 调度/执行对应的Router消息处理方法
	DoMsgHandle(request XRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router XRouter)
}
