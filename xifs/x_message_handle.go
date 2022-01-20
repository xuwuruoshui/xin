package xifs

// 消息管理
type XMessageHandle interface {

	// 调度/执行对应的Router消息处理方法
	DoMsgHandle(request XRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router XRouter)

	//启动worker工作池
	StartWorkerPool()

	// 将消息交给消息队列处理,由Worker进行处理
	SendMsgToTaskQueue(request XRequest)
}
